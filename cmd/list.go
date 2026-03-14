package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/oinkywaddles/earnings-cal-cli/internal/finnhub"
	"github.com/oinkywaddles/earnings-cal-cli/internal/index"
	"github.com/oinkywaddles/earnings-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	flagToday    bool
	flagTomorrow bool
	flagThisWeek bool
	flagNextWeek bool
	flagThisMonth bool
	flagNextMonth bool
	flagFrom     string
	flagTo       string
	flagSymbols  string
	flagAll      bool
	flagHour     string
	flagLimit    int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List upcoming and recent earnings",
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolVar(&flagToday, "today", false, "Today only")
	listCmd.Flags().BoolVar(&flagTomorrow, "tomorrow", false, "Tomorrow only")
	listCmd.Flags().BoolVar(&flagThisWeek, "this-week", false, "This week (default)")
	listCmd.Flags().BoolVar(&flagNextWeek, "next-week", false, "Next week")
	listCmd.Flags().BoolVar(&flagThisMonth, "this-month", false, "This month")
	listCmd.Flags().BoolVar(&flagNextMonth, "next-month", false, "Next month")
	listCmd.Flags().StringVar(&flagFrom, "from", "", "Start date (YYYY-MM-DD)")
	listCmd.Flags().StringVar(&flagTo, "to", "", "End date (YYYY-MM-DD)")
	listCmd.Flags().StringVar(&flagSymbols, "symbols", "", "Filter by symbols (comma-separated)")
	listCmd.Flags().BoolVar(&flagAll, "all", false, "No index filter, show all stocks")
	listCmd.Flags().StringVar(&flagHour, "hour", "", "Filter by hour: bmo, amc, dmh")
	listCmd.Flags().IntVar(&flagLimit, "limit", 0, "Limit number of results")

	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	apiKey := requireAPIKey()
	client := finnhub.NewClient(apiKey)

	from, to := resolveTimeRange()

	earnings, err := client.EarningsCalendar(from, to, "")
	if err != nil {
		exitWithError(err.Error())
	}

	// Determine symbol filter
	symbolFilter := resolveSymbolFilter()

	// Apply filters
	earnings = filterEarnings(earnings, symbolFilter)

	// Sort by date ascending
	sort.Slice(earnings, func(i, j int) bool {
		return earnings[i].Date < earnings[j].Date
	})

	// Apply limit
	if flagLimit > 0 && len(earnings) > flagLimit {
		earnings = earnings[:flagLimit]
	}

	filterDesc := output.FilterDesc(parseSymbolsFlag(), flagAll, os.Getenv("EARNINGS_WATCHLIST"))

	if jsonOutput {
		output.PrintListJSON(os.Stdout, earnings, from, to, filterDesc)
	} else {
		output.PrintListTable(os.Stdout, earnings, from, to, filterDesc)
	}
}

func resolveTimeRange() (string, string) {
	now := time.Now()
	dateFmt := "2006-01-02"

	// Custom range takes priority
	if flagFrom != "" || flagTo != "" {
		if flagFrom == "" || flagTo == "" {
			exitWithError("both --from and --to are required for custom date range")
		}
		validateDateRange(flagFrom, flagTo)
		return flagFrom, flagTo
	}

	switch {
	case flagToday:
		d := now.Format(dateFmt)
		return d, d
	case flagTomorrow:
		d := now.AddDate(0, 0, 1).Format(dateFmt)
		return d, d
	case flagNextWeek:
		// Monday of next week
		daysUntilMonday := (8 - int(now.Weekday())) % 7
		if daysUntilMonday == 0 {
			daysUntilMonday = 7
		}
		monday := now.AddDate(0, 0, daysUntilMonday)
		friday := monday.AddDate(0, 0, 4)
		return monday.Format(dateFmt), friday.Format(dateFmt)
	case flagThisMonth:
		first := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		last := first.AddDate(0, 1, -1)
		return first.Format(dateFmt), last.Format(dateFmt)
	case flagNextMonth:
		first := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		last := first.AddDate(0, 1, -1)
		return first.Format(dateFmt), last.Format(dateFmt)
	default:
		// this-week (default): Monday to Friday of current week
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		monday := now.AddDate(0, 0, -(weekday - 1))
		friday := monday.AddDate(0, 0, 4)
		return monday.Format(dateFmt), friday.Format(dateFmt)
	}
}

func parseSymbolsFlag() []string {
	if flagSymbols == "" {
		return nil
	}
	parts := strings.Split(flagSymbols, ",")
	var result []string
	for _, p := range parts {
		s := strings.TrimSpace(strings.ToUpper(p))
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func resolveSymbolFilter() map[string]bool {
	// --all → no filter
	if flagAll {
		return nil
	}

	// --symbols flag
	if symbols := parseSymbolsFlag(); len(symbols) > 0 {
		m := make(map[string]bool)
		for _, s := range symbols {
			m[s] = true
		}
		return m
	}

	// EARNINGS_WATCHLIST env
	if wl := os.Getenv("EARNINGS_WATCHLIST"); wl != "" {
		m := make(map[string]bool)
		for _, s := range strings.Split(wl, ",") {
			s = strings.TrimSpace(strings.ToUpper(s))
			if s != "" {
				m[s] = true
			}
		}
		return m
	}

	// Default: index constituents (auto-fetch if missing/expired)
	status := index.GetCacheStatus()
	if !status.Exists || status.Expired {
		warn("Fetching index constituents from Wikipedia... (run 'earnings-cal-cli init' to pre-cache)")
	}
	constituents, err := index.GetConstituents(noCache)
	if err != nil {
		warn("Failed to fetch index constituents, showing unfiltered results: " + err.Error())
		return nil
	}

	m := make(map[string]bool)
	for _, s := range constituents.Merged {
		m[s] = true
	}
	return m
}

const maxDateRangeDays = 90

func validateDateRange(from, to string) {
	dateFmt := "2006-01-02"
	f, err1 := time.Parse(dateFmt, from)
	t, err2 := time.Parse(dateFmt, to)
	if err1 != nil {
		exitWithError("invalid --from date format, expected YYYY-MM-DD")
	}
	if err2 != nil {
		exitWithError("invalid --to date format, expected YYYY-MM-DD")
	}
	if t.Before(f) {
		exitWithError("--to must be after --from")
	}
	if t.Sub(f).Hours()/24 > maxDateRangeDays {
		exitWithError(fmt.Sprintf("date range exceeds %d days", maxDateRangeDays))
	}
}

func filterEarnings(earnings []finnhub.Earning, symbolFilter map[string]bool) []finnhub.Earning {
	var result []finnhub.Earning
	for _, e := range earnings {
		// Symbol filter
		if symbolFilter != nil && !symbolFilter[e.Symbol] {
			continue
		}
		// Hour filter
		if flagHour != "" && !strings.EqualFold(e.Hour, flagHour) {
			continue
		}
		result = append(result, e)
	}
	return result
}
