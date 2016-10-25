package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/paddyforan/go-cli-talk/proverbs"
)

func main() {
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

	// we can also read command line arguments
	var id string
	flag.StringVar(&id, "proverb", "", "the ID of the proverb to retrieve; leave empty for a random proverb")
	flag.Parse()

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
