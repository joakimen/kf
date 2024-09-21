package cmd

import (
	"fmt"

	"github.com/joakimen/kf/pkg/kf"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print configuration file path",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := kf.Config()
		if err != nil {
			fmt.Println("error reading configuration file:", err)
			return
		}
		fmt.Println(config)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
