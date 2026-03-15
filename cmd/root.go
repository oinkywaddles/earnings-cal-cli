package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version = "dev"

	jsonOutput bool
	noCache    bool
)

func getVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

var rootCmd = &cobra.Command{
	Use:     "earnings-cal-cli",
	Short:   "Earnings calendar CLI — track upcoming and past earnings reports",
	Version: getVersion(),
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "Force refresh cache")
}

func exitWithError(msg string) {
	if jsonOutput {
		fmt.Fprintf(os.Stdout, "{\"error\":%q}\n", msg)
	} else {
		fmt.Fprintln(os.Stderr, "Error: "+msg)
	}
	os.Exit(1)
}

func warn(msg string) {
	if !jsonOutput {
		fmt.Fprintln(os.Stderr, msg)
	}
}

func requireAPIKey() string {
	key := os.Getenv("FINNHUB_API_KEY")
	if key == "" {
		exitWithError("FINNHUB_API_KEY not set. Get a free key at https://finnhub.io/register")
	}
	return key
}
