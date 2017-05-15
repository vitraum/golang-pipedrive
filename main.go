package pipedrive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type endpoints struct {
	PipelineDeals string
	Deal          string
	Deals         string
	DealFilter    string
	Pipelines     string
	Stages        string
	Filters       string
}

type getEndpointFunc func(endpoint string) (*http.Response, error)
type putEndpointFunc func(endpoint string, data io.Reader) (*http.Response, error)

// API represents the information needed to access the Pipedrive API v1
type API struct {
	token       string
	Endpoints   endpoints
	getEndpoint getEndpointFunc
	putEndpoint putEndpointFunc
	logURL      func(url string)
}

// Option represents an option given to the API constructor
type Option func(*API) error

// NewAPI create a new API object from the given options
func NewAPI(options ...Option) (*API, error) {
	pd := &API{
		logURL: func(u string) {},
	}

	for _, option := range options {
		err := option(pd)
		if err != nil {
			return nil, err
		}
	}

	return pd, nil
}

// GenericResponse lets the user of this package decode the data themselves
type GenericResponse struct {
	apiResult
	Data json.RawMessage
}

// Urler generates an url given an offset.
type Urler func(offset int) (string, error)

// FetchGeneric calls the Pipedrive API with a GET request.
func (pd *API) FetchGeneric(urlGenerator Urler, results chan GenericResponse) error {
	offset := 0
	for {
		url, err := urlGenerator(offset)
		if err != nil {
			return err
		}
		if url == "" {
			return nil
		}

		res, err := pd.getEndpoint(url)
		if err != nil {
			return err
		}

		var response GenericResponse
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return err
		}
		results <- response

		if !response.AdditionalData.Pagination.MoreItemsInCollection {
			return err
		}

		offset += response.AdditionalData.Pagination.Limit
	}
}

// PutGeneric makes a PUT request to the Pipedrive API using the supplied endpoint and data.
func (pd *API) PutGeneric(endpoint string, data io.Reader, results chan GenericResponse) error {
	res, err := pd.putEndpoint(endpoint, data)
	if err != nil {
		return err
	}

	var response GenericResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}
	results <- response

	if response.AdditionalData.Pagination.MoreItemsInCollection {
		return errors.New("Don't know how to handle MoreItemsInCollection after PUT")
	}
	return nil
}

// FetchDeals returns a list of deals, optionally using a filter
func (pd *API) FetchDeals(filterID int) (DealRefs, error) {
	var deals DealRefs
	start := 0
	for {
		url := fmt.Sprintf(pd.Endpoints.DealFilter, start, filterID)
		res, err := pd.getEndpoint(url)
		if err != nil {
			return nil, err
		}

		var pres struct {
			apiResult
			Data DealRefs
		}
		err = json.NewDecoder(res.Body).Decode(&pres)
		if err != nil {
			return nil, err
		}
		deals = append(deals, pres.Data...)

		if !pres.AdditionalData.Pagination.MoreItemsInCollection {
			return deals, nil
		}

		start += pres.AdditionalData.Pagination.Limit
	}
}

// FetchDeal returns a list of deals, optionally using a filter
func (pd *API) FetchDeal(dealID int) (DealRef, error) {
	url := fmt.Sprintf(pd.Endpoints.Deal, dealID)
	res, err := pd.getEndpoint(url)
	if err != nil {
		return DealRef{}, err
	}

	var pres struct {
		apiResult
		Data DealRef
	}

	var buf bytes.Buffer
	tee := io.TeeReader(res.Body, &buf)
	err = json.NewDecoder(tee).Decode(&pres)

	if err != nil {
		fmt.Println(&buf)
		return DealRef{}, err
	}

	return pres.Data, nil
}

// FetchDealsFromPipeline returns a list of all deals in a pipeline, optionally using a filter
func (pd *API) FetchDealsFromPipeline(plID, filterID int) (Deals, error) {
	var deals Deals
	start := 0
	for {
		url := fmt.Sprintf(pd.Endpoints.PipelineDeals, plID, start, filterID)
		res, err := pd.getEndpoint(url)
		if err != nil {
			return nil, err
		}

		var pres struct {
			apiResult
			Data Deals
		}
		err = json.NewDecoder(res.Body).Decode(&pres)
		if err != nil {
			return nil, err
		}
		deals = append(deals, pres.Data...)

		if !pres.AdditionalData.Pagination.MoreItemsInCollection {
			return deals, nil
		}

		start += pres.AdditionalData.Pagination.Limit
	}
}

// FetchDealUpdates request updates for a specific deal
func (pd *API) FetchDealUpdates(dealID int) (DealUpdates, error) {
	var dealUpdates DealUpdates
	start := 0
	for {
		url := fmt.Sprintf(pd.Endpoints.Deals, dealID, start)
		res, err := pd.getEndpoint(url)
		if err != nil {
			return nil, err
		}

		var pres struct {
			apiResult
			Data DealUpdates
		}

		err = json.NewDecoder(res.Body).Decode(&pres)
		if err != nil {
			return nil, err
		}
		dealUpdates = append(dealUpdates, pres.Data...)

		if !pres.AdditionalData.Pagination.MoreItemsInCollection {
			return dealUpdates, nil
		}

		start += pres.AdditionalData.Pagination.Limit
	}
}

// FetchPipelineChanges generates a list of deals with changes between the given stages
func (pd *API) FetchPipelineChanges(deals []Deal, stages Stages) (PipelineChangeResults, error) {
	res := make([]PipelineChangeResult, 0, len(deals))
	for i, deal := range deals {
		dealFlow := PipelineChangeResult{
			Deal: deals[i],
		}
		item := DealFlowUpdate{
			PiT:   dealFlow.Deal.Added,
			Phase: stages[0].Name,
		}
		dealFlow.Updates = append(dealFlow.Updates, item)
		updates, err := pd.FetchDealUpdates(deal.ID)
		if err != nil {
			return nil, err
		}
		for _, update := range updates {
			if update.StoryData.ActionType != "edit" || len(update.StoryData.ChangeLog) == 0 {
				continue
			}
			change := update.StoryData.ChangeLog[0]
			if change.FieldName == "Phase" {
				item := DealFlowUpdate{
					PiT:   update.StoryData.AddTime,
					Phase: change.NewValue.(string),
				}
				dealFlow.Updates = append(dealFlow.Updates, item)
			}
		}
		if len(dealFlow.Updates) > 0 {
			res = append(res, dealFlow)
		}
	}
	return res, nil
}

// GetPipelineIDByName searches for the given pipeline and returns its ID.
func (pd *API) GetPipelineIDByName(name string) (int, error) {
	res, err := pd.getEndpoint(pd.Endpoints.Pipelines)
	if err != nil {
		return 0, err
	}

	var pres struct {
		apiResult
		Data []Pipeline `json:"data"`
	}

	err = json.NewDecoder(res.Body).Decode(&pres)
	if err != nil {
		return 0, err
	}

	for _, pl := range pres.Data {
		if pl.Name != name {
			continue
		}
		return pl.ID, nil
	}

	return 0, fmt.Errorf("Pipeline '%s' not found", name)
}

// RetrieveStagesForPipeline returns all stages for a given pipeline
func (pd *API) RetrieveStagesForPipeline(plID int) (Stages, error) {
	res, err := pd.getEndpoint(fmt.Sprintf(pd.Endpoints.Stages, plID))
	if err != nil {
		return nil, err
	}

	var sres StagesResult
	err = json.NewDecoder(res.Body).Decode(&sres)
	if err != nil {
		return nil, err
	}

	return sres.Data, nil
}

// GetFilterIDByName returns the filter id for the given name
func (pd *API) GetFilterIDByName(name string) (int, error) {
	res, err := pd.getEndpoint(pd.Endpoints.Filters)
	if err != nil {
		return 0, err
	}

	var pres struct {
		apiResult
		Data Filters `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&pres)
	if err != nil {
		return 0, err
	}

	for _, pl := range pres.Data {
		if pl.Name != name {
			continue
		}
		return pl.Id, nil
	}

	return 0, fmt.Errorf("Pipeline '%s' not found", name)
}

// GetDealFieldByID returns the DealField with the given ID.
func (pd *API) GetDealFieldByID(id int) (DealField, error) {
	res, err := pd.getEndpoint(pd.Endpoints.Filters)
	if err != nil {
		return DealField{}, err
	}

	var pres struct {
		apiResult
		Data DealField
	}
	err = json.NewDecoder(res.Body).Decode(&pres)
	if err != nil {
		return DealField{}, err
	}

	return pres.Data, nil
}
