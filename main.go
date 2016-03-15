package pipedrive

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type endpoints struct {
	pipelineDeals string
	deals         string
	pipelines     string
	stages        string
	filters       string
}

type API struct {
	token       string
	eps         endpoints
	getEndpoint func(endpoint string) (*http.Response, error)
	logURL      func(url string)
}

// Option represents an option given to the API constructor
type Option func(*API) error

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

func (pd *API) FetchDealsFromPipeline(plID, filterID int) (Deals, error) {
	var deals Deals
	start := 0
	for {
		url := fmt.Sprintf(pd.eps.pipelineDeals, plID, start, filterID)
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
	return deals, nil
}

func (pd *API) FetchDealUpdates(dealId int) (DealUpdates, error) {
	var dealUpdates DealUpdates
	start := 0
	for {
		url := fmt.Sprintf(pd.eps.deals, dealId, start)
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
	return dealUpdates, nil
}

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
		updates, err := pd.FetchDealUpdates(deal.Id)
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

func (pd *API) GetPipelineIdByName(name string) (int, error) {
	res, err := pd.getEndpoint(pd.eps.pipelines)
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
		return pl.Id, nil
	}

	return 0, fmt.Errorf("Pipeline '%s' not found", name)
}

func (pd *API) RetrieveStagesForPipeline(plID int) (Stages, error) {
	res, err := pd.getEndpoint(fmt.Sprintf(pd.eps.stages, plID))
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

func (pd *API) GetFilterIdByName(name string) (int, error) {
	res, err := pd.getEndpoint(pd.eps.filters)
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
