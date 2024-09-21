package cmd

import (
	"fmt"
	"os"

	"github.com/joakimen/kf/pkg/userconfig"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all known files",
	Run: func(cmd *cobra.Command, args []string) {
		lines, err := userconfig.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading configuration file: %v\n", err)
			os.Exit(1)
		}

		for _, line := range lines {
			fmt.Println(line)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
