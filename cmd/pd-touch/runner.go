package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	pipedrive "github.com/vitraum/golang-pipedrive"
	"golang.org/x/sync/errgroup"
)

var (
	verbose = false
)

func main() {
	var token = ""
	flag.StringVar(&token, "token", "", "API token to be used (mandatory)")

	flag.BoolVar(&verbose, "verbose", verbose, "enable verbose output")

	flag.Parse()

	apiOptions := []pipedrive.Option{
		pipedrive.HTTPFetcher,
		pipedrive.WithCustomOrgFields(),
		pipedrive.WithCustomDealFields(),
	}

	switch token {
	case "":
		apiOptions = append(apiOptions, pipedrive.EnvToken(""))
	default:
		apiOptions = append(apiOptions, pipedrive.FixedToken(token))
	}

	pd, err := pipedrive.NewAPI(apiOptions...)
	if err != nil {
		logrus.Fatal(err)
	}

	ctx := context.Background()
	//	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	out := make(chan *pipedrive.DealRef)

	g.Go(func() error {
		defer close(out)
		for _, dealString := range flag.Args() {
			dealID, err := strconv.Atoi(dealString)
			if err != nil {
				return err
			}
			if dealID == 0 {
				return errors.New("DealID 0 not allowed")
			}

			start := time.Now()
			deal, err := pd.FetchDeal(dealID)
			if err != nil {
				return err
			}
			out <- &deal
			time.Sleep(time.Until(start.Add(333 * time.Millisecond)))
		}
		return nil
	})
	g.Go(func() error {
		for d := range out {
			aktion := d.CustomFields["Aktion"]

			newA := map[string]struct{}{}
			switch aktion {
			case nil:
				newA["839"] = struct{}{}
				err := setAktionen(fmt.Sprintf("%d", d.ID), newA)
				if err != nil {
					return err
				}
			default:
				newA = str2set(aktion)
				if _, e := newA["839"]; e {
					delete(newA, "839")
				} else {
					newA["839"] = struct{}{}
				}
				err = setAktionen(fmt.Sprintf("%d", d.ID), newA)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	err = g.Wait()
	if err != nil {
		logrus.Fatal(err)
	}
}

var client = http.Client{
	Timeout: 5 * time.Second,
}

func str2set(val interface{}) map[string]struct{} {
	out := make(map[string]struct{})
	csv, ok := val.(string)
	if !ok {
		return out
	}
	vals := strings.Split(csv, ",")
	for _, v := range vals {
		out[v] = struct{}{}
	}
	return out
}

func setAktionen(id string, aktionen map[string]struct{}) error {
	type Payload struct {
		Aktion string `json:"7de67a2875cf1fee9aa92dd0f8c65f5b24226b34"`
	}

	as := make([]string, 0, len(aktionen))
	for k := range aktionen {
		as = append(as, k)
	}

	data := Payload{
		Aktion: strings.Join(as, ","),
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)
	if verbose {
		logrus.Infof("updating %s to %s", id, as)
	}
	url := os.ExpandEnv(fmt.Sprintf("https://api.pipedrive.com/v1/deals/%v?api_token=$PDTOKEN", id))
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
