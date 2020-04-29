package pipedrive

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

var ErrEmptyToken = errors.New("Pipedrive token must not be empty")

var defaultEndpoints = endpoints{
	PipelineDeals:  "https://api.pipedrive.com/v1/pipelines/%d/deals?everyone=0&start=%d&filter_id=%d&limit=500",
	Deals:          "https://api.pipedrive.com/v1/deals/%d/updates?start=%d",
	Deal:           "https://api.pipedrive.com/v1/deals/%d",
	DealFilter:     "https://api.pipedrive.com/v1/deals?start=%d&filter_id=%d&limit=500",
	DealActivities: "https://api.pipedrive.com/v1/deals/%d/activities?start=%d",
	Pipelines:      "https://api.pipedrive.com/v1/pipelines",
	Stages:         "https://api.pipedrive.com/v1/stages?pipeline_id=%d",
	Filters:        "https://api.pipedrive.com/v1/filters",
	DealFields:     "https://api.pipedrive.com/v1/dealFields",
	DealField:      "https://api.pipedrive.com/v1/dealFields/%d",
	Organization:   "https://api.pipedrive.com/v1/organizations/%d",
}

func LogURLs(a *API) error {
	r := regexp.MustCompile(fmt.Sprintf("api_token=%s", a.token))
	a.logURL = func(u string) {
		prettyURL := r.ReplaceAllString(u, "api_token=…")
		fmt.Println(prettyURL)
	}

	return nil
}

func CustomURLLogger(logger func(u string), elipsifyToken bool) Option {
	return func(a *API) error {
		r := regexp.MustCompile(fmt.Sprintf("api_token=%s", a.token))
		a.logURL = func(u string) {
			prettyURL := u
			if elipsifyToken {
				prettyURL = r.ReplaceAllString(u, "api_token=…")
			}
			logger(prettyURL)
		}

		return nil
	}
}

func FixedToken(token string) Option {
	return func(a *API) error {
		if token == "" {
			return ErrEmptyToken
		}
		a.token = token
		return nil
	}
}

func EnvToken(envName string) Option {
	return func(a *API) error {
		if envName == "" {
			envName = "PDTOKEN"
		}

		token := os.Getenv(envName)
		if token == "" {
			return ErrEmptyToken
		}
		a.token = token

		return nil
	}
}

func (a *API) endpointFuncWithClient(get getEndpointFunc) getEndpointFunc {
	return func(endpoint string) (*http.Response, error) {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		values := u.Query()
		values.Add("api_token", a.token)
		u.RawQuery = values.Encode()
		a.logURL(u.String())
		res, err := get(u.String())
		return res, err
	}
}

func (a *API) requestEndpointFuncWithClient(doer func(req *http.Request) (*http.Response, error), method string) putEndpointFunc {
	return func(endpoint string, data io.Reader) (*http.Response, error) {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		values := u.Query()
		values.Add("api_token", a.token)
		u.RawQuery = values.Encode()
		a.logURL(u.String())
		req, err := http.NewRequest(method, u.String(), data)
		if err != nil {
			return nil, err
		}
		return doer(req)
	}
}

func HTTPFetcher(a *API) error {
	a.Endpoints = defaultEndpoints

	a.getEndpoint = a.endpointFuncWithClient(http.Get)

	client := http.Client{}
	a.putEndpoint = a.requestEndpointFuncWithClient(client.Do, "PUT")
	return nil
}

func HTTPFetcherWithTimeout(timeout time.Duration) Option {
	return func(a *API) error {
		a.Endpoints = defaultEndpoints

		client := http.Client{
			Timeout: timeout,
		}
		a.getEndpoint = a.endpointFuncWithClient(client.Get)
		a.putEndpoint = a.requestEndpointFuncWithClient(client.Do, "PUT")
		return nil
	}
}

func WithCustomOrgFields() Option {
	fields := make(map[string]cFieldExtractor)

	initOrgFields := func(a *API) error {
		fields["460a5e9346f7b7bf008904a285414900bd70ecbc"] = func(v interface{}) (string, interface{}, error) {
			if v == nil {
				return "E-Mail", "", nil
			}
			if vv, ok := v.(string); ok {
				return "E-Mail", vv, nil
			}
			return "", nil, fmt.Errorf("could not map: %#v", v)
		}
		return nil
	}

	mapOrgFields := func(o *Organization, jv map[string]interface{}) {
		if o.CustomFields == nil {
			o.CustomFields = make(map[string]interface{})
		}
		for key, ex := range fields {
			v, e := jv[key]
			if !e {
				continue
			}
			k, v, err := ex(jv[key])
			if err != nil {
				log.Print(err)
			}
			o.CustomFields[k] = v
		}
	}

	return func(a *API) error {
		a.afterInit = append(a.afterInit, initOrgFields)
		a.mapFieldsOrg = mapOrgFields
		return nil
	}
}

type cFieldExtractor func(v interface{}) (string, interface{}, error)

func WithCustomDealFields() Option {
	fields := make(map[string]cFieldExtractor)

	initDealFields := func(a *API) error {
		fields["66dd77c15a4867a95be45bc3ecc162fdf74a4c76"] = func(v interface{}) (string, interface{}, error) {
			if v == nil {
				return "Kunde to Go", "", nil
			}
			if vv, ok := v.(string); ok {
				switch vv {
				case "758":
					return "Kunde to Go", "Ja", nil
				case "759":
					return "Kunde to Go", "Nein", nil
				}
			}
			return "", nil, fmt.Errorf("could not map: %#v", v)
		}
		fields["7de67a2875cf1fee9aa92dd0f8c65f5b24226b34"] = func(v interface{}) (string, interface{}, error) {
			if v == nil {
				return "Aktion", "", nil
			}
			if vv, ok := v.(string); ok {
				return "Aktion", vv, nil
			}
			return "", nil, fmt.Errorf("could not map: %#v", v)
		}
		fields["1ed188d19ec50c6563dbdad533126dc58882b429"] = func(v interface{}) (string, interface{}, error) {
			if v == nil {
				return "Lead - Quelle / Medium", "", nil
			}
			if vv, ok := v.(string); ok {
				return "Lead - Quelle / Medium", vv, nil
			}
			return "", nil, fmt.Errorf("could not map: %#v", v)
		}

		return nil
	}

	mapDealFields := func(d *DealRef, jv map[string]interface{}) {
		if d.CustomFields == nil {
			d.CustomFields = make(map[string]interface{})
		}
		for key, ex := range fields {
			v, e := jv[key]
			if !e {
				continue
			}
			k, v, err := ex(jv[key])
			if err != nil {
				log.Print(err)
			}
			d.CustomFields[k] = v
		}
	}

	return func(a *API) error {
		a.afterInit = append(a.afterInit, initDealFields)
		a.mapFieldsDeal = mapDealFields
		return nil
	}
}
