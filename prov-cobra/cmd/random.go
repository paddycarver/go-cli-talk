package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Retrieve a random proverb.",
	Long:  "Retrieve a random proverb.",
	RunE: func(cmd *cobra.Command, args []string) error {
		proverb, err := getProverb("")
		if err != nil {
			return err
		}
		output := proverb.Value
		if verbose {
			output = fmt.Sprintf("%s: %s", proverb.ID, output)
		}
		fmt.Println(output)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(randomCmd)
}
