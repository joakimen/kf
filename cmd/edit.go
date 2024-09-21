package cmd

import (
	"fmt"
	"os"

	"github.com/joakimen/kf/internal/config"
	"github.com/joakimen/kf/internal/fs"
	"github.com/joakimen/kf/internal/fuzzy"
	"github.com/joakimen/kf/internal/slice"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a known file",
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

		selectedFile := fuzzy.SelectFile(sanitizedConfigLines)
		fmt.Println(selectedFile)
		err = fs.EditFile(fs.GetEditorName(), selectedFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error editing file: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
