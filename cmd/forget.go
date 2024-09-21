package cmd

import (
	"fmt"
	"os"

	"github.com/joakimen/kf/pkg/config"

	"github.com/spf13/cobra"
)

var forgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget a file from the list of known files",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		fileToRemove := args[0]
		removedLines, err := config.RemoveEntry(fileToRemove)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error removing from configuration file: %v\n", err)
			os.Exit(1)
		}

		if len(removedLines) > 0 {
			fmt.Println("Removed the following entries:")
			for _, line := range removedLines {
				fmt.Printf("%3d: %s\n", line.Number, line.Text)
			}
		} else {
			fmt.Println("No matching entries found:", fileToRemove)
		}
	},
}

func init() {
	rootCmd.AddCommand(forgetCmd)
}
