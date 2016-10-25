package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/paddyforan/go-cli-talk/proverbs"
)

func printUsage() {
	yellow := color.New(color.FgYellow)
	printYellow := yellow.SprintFunc()
	fmt.Printf("%s prov-stdlib [-v] COMMAND [ID]\n\nSupported commands:\n  random\tReturn a random proverb.\n  get\t\tReturn the proverb associated with the passed ID.\n  help\t\tPrint this message.\n\nSupported flags:\n", printYellow("Usage:"))
	flag.PrintDefaults()
}

func main() {
	// let's have a verbose mode
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Print the proverb's ID as well as the proverb.")
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
	if action == "help" {
		printUsage()
		return
	}
	if action != "help" && action != "get" && action != "random" {
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
		fmt.Printf("%s PROVERBS_URL must be set to the API endpoint to retrieve proverbs from.\n", printBoldRed("[ERROR]"))
		os.Exit(1) // exit codes are a simple call
	}

	headers := http.Header{}

	errMode := os.Getenv("ERROR")

	switch errMode {
	case "400":
		headers.Set("Return-Error", "bad-request")
	case "500":
		headers.Set("Return-Error", "internal")
	}

	proverb, err := proverbs.GetProverb(baseURL, id, headers)
	if err != nil {
		fmt.Printf("%s %s\n", printBoldRed("[ERROR]"), err.Error())
		os.Exit(1)
	}
	output := proverb.Value
	if verbose {
		output = fmt.Sprintf("%s: %s", proverb.ID, output)
	}
	fmt.Println(output)
}
