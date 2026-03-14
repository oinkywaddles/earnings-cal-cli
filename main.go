package main

import (
	"os"

	"github.com/oinkywaddles/earnings-cal-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
