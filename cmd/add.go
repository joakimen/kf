package cmd

import (
	"fmt"
	"github.com/joakimen/kf/internal/config"
	"github.com/joakimen/kf/internal/fs"
	"github.com/spf13/cobra"
	"os"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a file to the list of known files",
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
		Run: func(cmd *cobra.Command, args []string) {
			fileToAddAbs, err := fs.RealPath(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting real path: %v\n", err)
				os.Exit(1)
			}

			if !fs.IsValidFile(fileToAddAbs) {
				fmt.Fprintf(os.Stderr, "not a valid file: %s\n", fileToAddAbs)
				os.Exit(1)
			}

			err = config.AppendToConfigFile(fileToAddAbs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error appending to configuration file: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("kf: added", args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
