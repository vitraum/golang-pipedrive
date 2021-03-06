package pipedrive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type endpoints struct {
	PipelineDeals  string
	Deal           string
	Deals          string
	DealFilter     string
	Pipelines      string
	Stages         string
	Filters        string
	DealField      string
	DealActivities string
	Organization   string

	DealFields         string
	OrganizationFields string
}

type getEndpointFunc func(endpoint string) (*http.Response, error)
type putEndpointFunc func(endpoint string, data io.Reader) (*http.Response, error)

// API represents the information needed to access the Pipedrive API v1
type API struct {
	token       string
	Endpoints   endpoints
	getEndpoint getEndpointFunc
	putEndpoint putEndpointFunc

	afterInit []Option

	mapFieldsOrg  func(*Organization, map[string]interface{})
	mapFieldsDeal func(*DealRef, map[string]interface{})

	logURL func(url string)
}

// Option represents an option given to the API constructor
type Option func(*API) error

// NewAPI create a new API object from the given options
func NewAPI(options ...Option) (*API, error) {
	pd := &API{
		logURL:    func(u string) {},
		afterInit: make([]Option, 0),
	}

	for _, option := range options {
		err := option(pd)
		if err != nil {
			return nil, err
		}
	}

	for _, initF := range pd.afterInit {
		err := initF(pd)
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
		err2 := res.Body.Close()
		if err != nil {
			return err
		}
		if err2 != nil {
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
			Data []json.RawMessage
		}
		err = json.NewDecoder(res.Body).Decode(&pres)
		if err != nil {
			return nil, err
		}
		for _, data := range pres.Data {
			var dr DealRef
			err = json.Unmarshal(data, &dr)
			if err != nil {
				return nil, err
			}
			if pd.mapFieldsDeal != nil {
				var jv map[string]interface{}
				err = json.Unmarshal(data, &jv)
				if err != nil {
					return nil, err
				}
				pd.mapFieldsDeal(&dr, jv)
			}
			deals = append(deals, dr)
		}

		if !pres.AdditionalData.Pagination.MoreItemsInCollection {
			return deals, nil
		}

		start += pres.AdditionalData.Pagination.Limit
	}
}

// FetchDeal returns a deal with the given id
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
		logrus.Errorf("Error decoding result: %s", buf.String())
		return DealRef{}, err
	}

	dr := pres.Data
	if pd.mapFieldsDeal != nil {
		mf := struct {
			Data map[string]interface{}
		}{}
		err = json.Unmarshal(buf.Bytes(), &mf)
		if err != nil {
			return DealRef{}, err
		}
		pd.mapFieldsDeal(&dr, mf.Data)
	}

	return dr, nil
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

	in := make(chan Deal)
	out := make(chan PipelineChangeResult)
	errors := make(chan error)

	wg := sync.WaitGroup{}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for deal := range in {
				dealFlow := PipelineChangeResult{
					Deal: deal,
				}
				item := DealFlowUpdate{
					PiT:   dealFlow.Deal.Added,
					Phase: stages[0].Name,
				}
				dealFlow.PipelineUpdates = append(dealFlow.PipelineUpdates, item)
				var err error
				dealFlow.Updates, err = pd.FetchDealUpdates(deal.ID)
				if err != nil {
					errors <- err
				}
				for _, update := range dealFlow.Updates {
					if update.StoryData.ActionType != "edit" || len(update.StoryData.ChangeLog) == 0 {
						continue
					}
					change := update.StoryData.ChangeLog[0]
					if change.FieldName == "Phase" {
						item := DealFlowUpdate{
							PiT:   update.StoryData.AddTime,
							Phase: change.NewValue.(string),
						}
						dealFlow.PipelineUpdates = append(dealFlow.PipelineUpdates, item)
					}
				}

				if len(dealFlow.PipelineUpdates) > 0 {
					out <- dealFlow
				}
			}
		}()
	}

	go func() {
		for _, deal := range deals {
			in <- deal
		}
		close(in)
	}()

	go func() {
		for df := range out {
			res = append(res, df)
		}
	}()

	errs := []error{}
	go func() {
		for err := range errors {
			errs = append(errs, err)
		}
	}()

	wg.Wait()
	close(out)
	close(errors)

	if len(errs) > 0 {
		return nil, errs[0]
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
	res, err := pd.getEndpoint(fmt.Sprintf(pd.Endpoints.DealField, id))
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

func (pd *API) GenericStreamHelper(worker func(r GenericResponse) error, generator Urler, closer func()) <-chan error {
	errs := make(chan error)
	wg := sync.WaitGroup{}

	// TODO add context for cancelation

	genresults := make(chan GenericResponse)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := pd.FetchGeneric(generator, genresults)
		close(genresults)
		if err != nil {
			errs <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range genresults {
			err := worker(r)
			if err != nil {
				errs <- err
			}
		}
		closer()
	}()

	go func() {
		wg.Wait()
		close(errs)
	}()

	return errs
}

func (pd *API) FetchDealActivities(dealID int) (<-chan Activity, <-chan error) {
	results := make(chan Activity)

	generator := func(offset int) (string, error) {
		return fmt.Sprintf(pd.Endpoints.DealActivities, dealID, offset), nil
	}
	worker := func(r GenericResponse) error {
		result := []Activity{}
		err := json.NewDecoder(bytes.NewReader(r.Data)).Decode(&result)
		if err != nil {
			return err
		}
		for _, a := range result {
			results <- a
		}
		return nil
	}
	closer := func() { close(results) }
	errs := pd.GenericStreamHelper(worker, generator, closer)

	return results, errs
}

// FetchOrganization returns the Organization with the given id
func (pd *API) FetchOrganization(orgID int) (Organization, error) {
	url := fmt.Sprintf(pd.Endpoints.Organization, orgID)
	res, err := pd.getEndpoint(url)
	if err != nil {
		return Organization{}, err
	}

	var pres struct {
		apiResult
		Data Organization
	}

	var buf bytes.Buffer
	tee := io.TeeReader(res.Body, &buf)
	err = json.NewDecoder(tee).Decode(&pres)
	if err != nil {
		logrus.Errorf("Error decoding result: %s", buf.String())
		return Organization{}, err
	}
	o := pres.Data

	if pd.mapFieldsOrg != nil {
		var cvres struct {
			apiResult
			Data map[string]interface{}
		}

		err = json.Unmarshal(buf.Bytes(), &cvres)
		if err != nil {
			logrus.Errorf("Error decoding result: %s", buf.String())
			return Organization{}, err
		}
		pd.mapFieldsOrg(&o, cvres.Data)
	}

	return o, nil
}
