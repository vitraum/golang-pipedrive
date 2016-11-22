package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path"
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

	if token == "" {
		execname := path.Base(os.Args[0])
		fmt.Printf("%s prints information about all deals\nUsage: %s [dealid] [dealid]...\n",
			execname, execname)
		flag.PrintDefaults()
		os.Exit(1)
	}

	apiOptions := []pipedrive.Option{
		pipedrive.HTTPFetcher,
		pipedrive.FixedToken(token),
	}

	if filterID > 0 && flag.NArg() > 0 {
		fmt.Println("Fatal error: filter and explicit dealIDs are mutually exclusive\n")
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
	tmpl = template.Must(template.New("name").Parse(outputTemplate))

	pd, err := pipedrive.NewAPI(apiOptions...)
	if err != nil {
		panic(err)
	}

	var deals []interface{}
	if flag.NArg() > 0 {
		deals, err = selectDeals(pd, flag.Args())
		if err != nil {
			panic(err)
		}
	} else {
		alldeals, err := pd.FetchDeals(filterID)
		if err != nil {
			panic(err)
		}

		if sample > 0 {
			for i := len(deals); i < sample; i++ {
				deals = append(deals, alldeals[rand.Intn(len(alldeals))])
			}
		} else {
			for i := 0; i < len(alldeals); i++ {
				deals = append(deals, alldeals[i])
			}
		}
	}

	if showVariables {
		fmt.Printf("%+v\n", deals[0])
		os.Exit(0)
	}

	for _, deal := range deals {
		err = printDeal(deal)
		if err != nil {
			panic(err)
		}
	}
}

func printDeal(deal interface{}) error {
	return tmpl.Execute(os.Stdout, deal)
}

func selectDeals(pd *pipedrive.API, dealIDs []string) ([]interface{}, error) {
	deals := make([]interface{}, 0, len(dealIDs))
	for _, dealString := range dealIDs {
		dealID, err := strconv.Atoi(dealString)
		if err != nil {
			return nil, err
		}

		deal, err := pd.FetchDeal(dealID)
		if err != nil {
			return nil, err
		}
		deals = append(deals, deal)
	}
	return deals, nil
}
