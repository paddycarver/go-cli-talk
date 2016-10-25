package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/paddyforan/go-cli-talk/proverbs"
)

func printUsage() {
	yellow := color.New(color.FgYellow)
	printYellow := yellow.SprintFunc()
	fmt.Printf("%s prov-stdlib COMMAND [ID]\n\nSupported commands:\n\trandom\t\tReturn a random proverb.\n\tget $ID\t\tReturn the proverb associated with $ID.\n\thelp\t\tPrint this message.\n", printYellow("Usage:"))
}

func main() {
	// let's use subcommands
	var action, id string
	if len(os.Args) < 2 || len(os.Args) > 3 {
		printUsage()
		os.Exit(1)
	}
	action = os.Args[1]
	if len(os.Args) > 2 {
		id = os.Args[2]
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
	fmt.Println(proverb.Value)
}
