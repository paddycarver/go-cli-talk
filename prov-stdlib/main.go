package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/paddyforan/go-cli-talk/proverbs"
)

func printUsage() {
	yellow := color.New(color.FgYellow)
	printYellow := yellow.SprintFunc()
	fmt.Fprintf(color.Output, "%s prov-stdlib [-v] COMMAND [ID]\n\nSupported commands:\n  random\tReturn a random proverb.\n  get\t\tReturn the proverb associated with the passed ID.\n  watch\t\tPrint a new proverb every second.\n  help\t\tPrint this message.\n\nSupported flags:\n", printYellow("Usage:"))
	flag.PrintDefaults()
}

func getProverb(baseURL, id, errMode, delay, chaos string, verbose bool, retries int) (proverbs.Quote, error) {
	headers := http.Header{}

	switch errMode {
	case "400":
		headers.Set("Return-Error", "bad-request")
	case "500":
		headers.Set("Return-Error", "internal")
	}

	if delay != "" {
		headers.Set("Sleep", delay)
	}

	if chaos != "" {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(99)%2 == 0 {
			headers.Set("Return-Error", "internal")
		}
	}

	var proverb proverbs.Quote
	var err error

	for attempt := 1; attempt <= retries; attempt++ {
		if verbose {
		}
		proverb, err = proverbs.GetProverb(baseURL, id, headers)
		if err != nil {
			if verbose {
				fmt.Fprintf(color.Output, "%s %s\n", "[ERROR]", err.Error())
			}
			if attempt == retries {
				return proverbs.Quote{}, err
			}
		} else {
			break
		}
	}
	return proverb, nil
}

func printProverb(proverb proverbs.Quote, verbose bool) {
	output := proverb.Value
	if verbose {
		output = fmt.Sprintf("%s: %s", proverb.ID, output)
	}
	fmt.Fprintln(color.Output, output)
}

func main() {
	// let's have a verbose mode
	var verbose bool
	var retry int
	flag.BoolVar(&verbose, "v", false, "Print the proverb's ID as well as the proverb.")
	flag.IntVar(&retry, "retry-attempts", 3, "Number of attempts to retry in the face of a server error.")
	flag.Parse()
	flag.Usage = printUsage

	// let's use subcommands
	var action, id string
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		printUsage()
		os.Exit(1)
	}
	action = args[0]
	if len(args) > 1 {
		id = args[1]
	}
	if action == "get" && id == "" {
		printUsage()
		os.Exit(1)
	}
	if action == "random" && id != "" {
		printUsage()
		os.Exit(1)
	}
	if action == "watch" && id != "" {
		printUsage()
		os.Exit(1)
	}
	if action == "help" {
		printUsage()
		return
	}
	if action != "help" && action != "get" && action != "random" && action != "watch" {
		printUsage()
		os.Exit(1)
	}

	// let's add some color
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)
	printBoldRed := boldRed.SprintFunc()

	// environment variables are simple to use
	baseURL := os.Getenv("PROVERBS_URL")
	if baseURL == "" {
		fmt.Fprintf(color.Output, "%s PROVERBS_URL must be set to the API endpoint to retrieve proverbs from.\n", printBoldRed("[ERROR]"))
		os.Exit(1) // exit codes are a simple call
	}
	errMode := os.Getenv("ERROR")
	delay := os.Getenv("HAMMERTIME")
	chaos := os.Getenv("CHAOS")

	if action == "watch" {
		for range time.Tick(time.Second) {
			proverb, err := getProverb(baseURL, id, errMode, delay, chaos, verbose, retry)
			if err != nil {
				fmt.Fprintf(color.Output, "%s %s\n", printBoldRed("[ERROR]"), err.Error())
				os.Exit(1)
			}

			printProverb(proverb, verbose)
		}
	} else {
		proverb, err := getProverb(baseURL, id, errMode, delay, chaos, verbose, retry)
		if err != nil {
			fmt.Fprintf(color.Output, "%s %s\n", printBoldRed("[ERROR]"), err.Error())
			os.Exit(1)
		}

		printProverb(proverb, verbose)
	}
}
