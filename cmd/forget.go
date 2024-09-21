package cmd

import (
	"fmt"
	"os"

	"github.com/joakimen/kf/pkg/kf"

	"github.com/spf13/cobra"
)

var forgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget a file from the list of known files",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		fileToRemove := args[0]
		matchingLineRemoved, err := kf.Forget(fileToRemove)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error removing from configuration file: %v\n", err)
			os.Exit(1)
		}

		if matchingLineRemoved {
			fmt.Println("Removed", fileToRemove)
			return
		}

		fmt.Println("No matching entries found. See 'list' command for existing entries.")
	},
}

func init() {
	rootCmd.AddCommand(forgetCmd)
}
