package cmd

import (
	"fmt"
	"os"

	"github.com/joakimen/kf/pkg/config"
	"github.com/joakimen/kf/pkg/slice"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all known files",
	Run: func(cmd *cobra.Command, args []string) {
		configFileLines, err := config.ReadConfigFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading configuration file: %v\n", err)
			os.Exit(1)
		}

		sanitizedConfigLines, err := slice.SanitizeFileSlice(configFileLines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sanitizing configuration lines: %v\n", err)
			os.Exit(1)
		}

		for _, line := range sanitizedConfigLines {
			fmt.Println(line)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
