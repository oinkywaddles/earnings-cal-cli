package cmd

import (
	"fmt"
	"time"

	"github.com/oinkywaddles/earnings-cal-cli/internal/index"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize or refresh index constituents cache",
	Long:  "Fetch S&P 100, NASDAQ 100, and Dow Jones constituents from Wikipedia and cache locally.",
	Run:   runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	status := index.GetCacheStatus()
	if status.Exists && !status.Expired && !noCache {
		fmt.Printf("Cache fresh: %d symbols, updated %s (use --no-cache to force)\n",
			status.Count, status.UpdatedAt.Format(time.DateOnly))
		return
	}

	c, err := index.GetConstituents(true)
	if err != nil {
		exitWithError("failed to fetch index constituents: " + err.Error())
	}

	fmt.Printf("Cached %d symbols (S&P 100: %d, NASDAQ 100: %d, Dow Jones: %d), valid 90 days\n",
		len(c.Merged), len(c.SP100), len(c.Nasdaq100), len(c.DowJones))
}
