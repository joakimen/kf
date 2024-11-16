package main

import (
	"os"

	"github.com/joakimen/kf/cmd"
)

func main() {
	run(os.Args, os.Getenv)
}

func run(args []string, getenv func(string) string) {
	app := cmd.NewApp(getenv)

	// pass getenv function to the app
	err := app.Run(args)
	if err != nil {
		os.Exit(1)
	}
}
