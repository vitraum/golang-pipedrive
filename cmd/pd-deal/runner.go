package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/vitraum/golang-pipedrive"
)

var (
	tmpl *template.Template
)

func main() {
	var token = ""
	flag.StringVar(&token, "token", "", "API token to be used (mandatory)")

	var verbose = false
	flag.BoolVar(&verbose, "verbose", verbose, "enable verbose output")

	var outputTemplate = "{{.ID}} {{.Status}}\n"
	flag.StringVar(&outputTemplate, "template", "", "text/template to be printed for each deal")

	var skipNewline = false
	flag.BoolVar(&skipNewline, "newline", skipNewline, "do not append newline to template")

	var sample = 0
	flag.IntVar(&sample, "sample", sample, "number of random samples to take")

	var seed int64
	flag.Int64Var(&seed, "seed", 0, "to be used for random sampling")

	var filterID = 0
	flag.IntVar(&filterID, "filter", 0, "filter ID to use for all deals")

	flag.Parse()

	if token == "" {
		fmt.Println("token is mandatory")
		flag.Usage()
		os.Exit(1)
	}

	apiOptions := []pipedrive.Option{
		pipedrive.HTTPFetcher,
		pipedrive.FixedToken(token),
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
	tmpl = template.Must(template.New("name").Parse(outputTemplate))

	pd, err := pipedrive.NewAPI(apiOptions...)
	if err != nil {
		panic(err)
	}

	if flag.NArg() > 0 {
		selectDeals(pd, flag.Args())
	} else {
		var deals pipedrive.DealRefs
		alldeals, err := pd.FetchDeals(filterID)
		if err != nil {
			panic(err)
		}

		if sample > 0 {
			for i := len(deals); i < sample; i++ {
				deals = append(deals, alldeals[rand.Intn(len(alldeals))])
			}
		} else {
			deals = alldeals
		}

		for _, deal := range deals {
			err = printDeal(deal)
			if err != nil {
				panic(err)
			}
		}
	}
}

func printDeal(deal interface{}) error {
	return tmpl.Execute(os.Stdout, deal)
}

func selectDeals(pd *pipedrive.API, deals []string) error {
	for _, dealString := range deals {
		dealID, err := strconv.Atoi(dealString)
		if err != nil {
			return err
		}

		deal, err := pd.FetchDeal(dealID)
		if err != nil {
			return err
		}

		err = printDeal(deal)
		if err != nil {
			return err
		}
	}
	return nil
}
