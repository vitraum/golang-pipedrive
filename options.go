package pipedrive

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

var ErrEmptyToken = errors.New("Token must not be empty")

var defaultEndpoints = endpoints{
	pipelineDeals: "https://api.pipedrive.com/v1/pipelines/%d/deals?everyone=0&start=%d&filter_id=%d",
	deals:         "https://api.pipedrive.com/v1/deals/%d/updates?start=%d",
	deal:          "https://api.pipedrive.com/v1/deals/%d",
	dealFilter:    "https://api.pipedrive.com/v1/deals?start=%d&filter_id=%d",
	pipelines:     "https://api.pipedrive.com/v1/pipelines",
	stages:        "https://api.pipedrive.com/v1/stages?pipeline_id=%d",
	filters:       "https://api.pipedrive.com/v1/filters",
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

func HTTPFetcher(a *API) error {
	a.eps = defaultEndpoints

	a.getEndpoint = a.endpointFuncWithClient(http.Get)
	return nil
}

func HTTPFetcherWithTimeout(timeout time.Duration) Option {
	return func(a *API) error {
		a.eps = defaultEndpoints

		client := http.Client{
			Timeout: timeout,
		}
		a.getEndpoint = a.endpointFuncWithClient(client.Get)
		return nil
	}
}
