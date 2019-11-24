package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"

	pipedrive "github.com/vitraum/golang-pipedrive"
)

var (
	tmpl *template.Template
)

func main() {
	var token = ""
	flag.StringVar(&token, "token", "", "API token to be used (mandatory)")

	var verbose = false
	flag.BoolVar(&verbose, "verbose", verbose, "enable verbose output")

	var outputTemplate = "{{.Id}} {{.Status}}"
	flag.StringVar(&outputTemplate, "template", outputTemplate, "text/template to be printed for each deal")

	var skipNewline = false
	flag.BoolVar(&skipNewline, "newline", skipNewline, "do not append newline to template")

	var sample = 0
	flag.IntVar(&sample, "sample", sample, "number of random samples to take")

	var seed int64
	flag.Int64Var(&seed, "seed", 0, "to be used for random sampling")

	var filterID = 0
	flag.IntVar(&filterID, "filter", 0, "filter ID to use for all deals")

	var showVariables = false
	flag.BoolVar(&showVariables, "showVariables", showVariables, "dump a sample deal")

	flag.Parse()

	apiOptions := []pipedrive.Option{
		pipedrive.HTTPFetcher,
	}

	switch token {
	case "":
		apiOptions = append(apiOptions, pipedrive.EnvToken(""))
	default:
		apiOptions = append(apiOptions, pipedrive.FixedToken(token))
	}

	if filterID > 0 && flag.NArg() > 0 {
		fmt.Println("Fatal error: filter and explicit dealIDs are mutually exclusive")
		os.Exit(1)
	}

	if showVariables {
		sample = 1
	}

	if sample > 0 {
		if seed == 0 {
			seed = time.Now().UTC().UnixNano()
		}
		if verbose {
			fmt.Printf("Using seed %d\n", seed)
		}
		rand.Seed(seed)
	}

	if verbose {
		apiOptions = append(apiOptions, pipedrive.LogURLs)
	}

	if verbose {
		fmt.Printf("using '%s' as template\n", outputTemplate)
	}

	if !skipNewline {
		outputTemplate = fmt.Sprintf("%s\n", outputTemplate)
	}

	pd, err := pipedrive.NewAPI(apiOptions...)
	if err != nil {
		panic(err)
	}

	tmpl = template.New("name")
	tmpl.Funcs(template.FuncMap{
		"Age": func(t pipedrive.Time) int { return int(time.Now().Sub(t.Time).Hours()) / 24 },
		"Org": func(id int) pipedrive.Organization {
			org, err := pd.FetchOrganization(id)
			if err != nil {
				panic(err)
			}
			return org
		},
	})
	tmpl = template.Must(tmpl.Parse(outputTemplate))

	deals := make(chan interface{})
	var wg sync.WaitGroup
	if flag.NArg() > 0 {
		wg.Add(1)

		fetchDeals := func(deals []string, out chan<- interface{}, copydone, done func()) {
			tmp := make([]string, len(deals))
			copy(tmp, deals)
			copydone()
			//logrus.Infof("selectDeals with %d deals", len(deals))
			err := selectDeals(pd, tmp, out)
			if err != nil {
				panic(err)
			}
			done()
		}

		go func() {
			defer wg.Done()
			cutoff := flag.NArg() / 5
			//logrus.Infof("cutoff %d %d", cutoff, flag.NArg())
			tmpDeals := make([]string, 0, cutoff)
			for _, dealID := range flag.Args() {
				tmpDeals = append(tmpDeals, dealID)
				if len(tmpDeals) <= cutoff {
					continue
				}
				var wgtmp sync.WaitGroup
				wg.Add(1)
				wgtmp.Add(1)
				go fetchDeals(tmpDeals, deals, func() { wgtmp.Done() }, func() { wg.Done() })
				wgtmp.Wait()
				tmpDeals = make([]string, 0, cutoff)
			}
			if len(tmpDeals) > 0 {
				wg.Add(1)
				fetchDeals(tmpDeals, deals, func() {}, func() { wg.Done() })
			}
		}()

	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			alldeals, err := pd.FetchDeals(filterID)
			if err != nil {
				panic(err)
			}

			if sample > 0 {
				for i := 0; i < sample; i++ {
					deals <- alldeals[rand.Intn(len(alldeals))]
				}
			} else {
				for _, d := range alldeals {
					deals <- d
				}
			}
		}()
	}

	if showVariables {
		fmt.Printf("%+v\n", <-deals)
		os.Exit(0)
	}

	go func() {
		wg.Wait()
		close(deals)
	}()

	for deal := range deals {
		err = printDeal(deal)
		if err != nil {
			panic(err)
		}
	}
}

func printDeal(deal interface{}) error {
	dj, err := json.Marshal(deal)
	if err != nil {
		panic(err)
	}

	jsonDeal := map[string]interface{}{}
	err = json.Unmarshal(dj, &jsonDeal)
	if err != nil {
		panic(err)
	}

	return tmpl.Execute(os.Stdout, deal)
}

func selectDeals(pd *pipedrive.API, dealIDs []string, out chan<- interface{}) error {
	for _, dealString := range dealIDs {
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
		out <- deal
		time.Sleep(time.Until(start.Add(1 * time.Second)))
	}
	return nil
}
