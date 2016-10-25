package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var verbose bool
var printBoldRed = color.New(color.FgRed).Add(color.Bold).SprintFunc()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "prov-cobra",
	Short:         "Fetch proverbs about Go from an API.",
	Long:          `An easy way to access the collected pithy wisdom that the Go community has accumulated.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if err == errMissingID {
			getCmd.Usage()
		} else {
			fmt.Fprintf(color.Output, "%s %s\n", printBoldRed("[ERROR]"), err.Error())
		}
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print the ID of a proverb as well as its text.")
}
