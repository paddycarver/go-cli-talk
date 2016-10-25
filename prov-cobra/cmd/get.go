package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var errMissingID = errors.New("ID must be specified")

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get {ID}",
	Short: "Retrieve a proverb.",
	Long:  "Retrieve the proverb specified by {ID}.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errMissingID
		}

		proverb, err := getProverb(args[0])
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
	RootCmd.AddCommand(getCmd)
}
