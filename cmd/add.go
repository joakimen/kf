package cmd

import (
	"fmt"
	"os"

	"github.com/joakimen/kf/pkg/kf"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a file to the list of known files",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		fileToAdd := args[0]
		err := kf.Add(fileToAdd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error adding to configuration file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Added", args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
