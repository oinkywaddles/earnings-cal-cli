package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/oinkywaddles/earnings-cal-cli/internal/finnhub"
	"github.com/oinkywaddles/earnings-cal-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	detailFrom string
	detailTo   string
)

var detailCmd = &cobra.Command{
	Use:   "detail SYMBOL [SYMBOL...]",
	Short: "Show earnings detail for one or more symbols",
	Args:  cobra.MinimumNArgs(1),
	Run:   runDetail,
}

func init() {
	detailCmd.Flags().StringVar(&detailFrom, "from", "", "Start date (YYYY-MM-DD)")
	detailCmd.Flags().StringVar(&detailTo, "to", "", "End date (YYYY-MM-DD)")

	rootCmd.AddCommand(detailCmd)
}

func runDetail(cmd *cobra.Command, args []string) {
	apiKey := requireAPIKey()
	client := finnhub.NewClient(apiKey)

	// Normalize symbols
	var symbols []string
	seen := make(map[string]bool)
	for _, a := range args {
		s := strings.ToUpper(strings.TrimSpace(a))
		if s != "" && !seen[s] {
			seen[s] = true
			symbols = append(symbols, s)
		}
	}

	const maxSymbols = 20
	if len(symbols) > maxSymbols {
		exitWithError(fmt.Sprintf("too many symbols (%d), max is %d", len(symbols), maxSymbols))
	}

	from, to := resolveDetailRange()

	// Fetch per symbol (server-side filtering, minimal data)
	var allEarnings []finnhub.Earning
	for _, sym := range symbols {
		earnings, err := client.EarningsCalendar(from, to, sym)
		if err != nil {
			exitWithError(err.Error())
		}
		allEarnings = append(allEarnings, earnings...)
	}

	// Pick most relevant: last past + next upcoming per symbol
	matched := findDetailEarnings(allEarnings, seen)

	if len(matched) == 0 {
		exitWithError("no earnings found for specified symbols in date range")
	}

	if jsonOutput {
		output.PrintDetailJSON(os.Stdout, matched)
	} else {
		output.PrintDetailText(os.Stdout, matched, nil)
	}
}

func resolveDetailRange() (string, string) {
	if detailFrom != "" || detailTo != "" {
		if detailFrom == "" || detailTo == "" {
			exitWithError("both --from and --to are required for custom date range")
		}
		validateDateRange(detailFrom, detailTo)
		return detailFrom, detailTo
	}

	// Default: ±3 months (covers last + next earnings per symbol)
	// Not subject to maxDateRangeDays since detail queries per-symbol
	now := time.Now()
	dateFmt := "2006-01-02"
	from := now.AddDate(0, -3, 0).Format(dateFmt)
	to := now.AddDate(0, 3, 0).Format(dateFmt)
	return from, to
}

// findDetailEarnings picks the most relevant earnings for each symbol:
// the most recent past + the next upcoming.
func findDetailEarnings(earnings []finnhub.Earning, symbols map[string]bool) []finnhub.Earning {
	today := time.Now().Format("2006-01-02")

	type pair struct {
		past   *finnhub.Earning
		future *finnhub.Earning
	}
	bySymbol := make(map[string]*pair)

	for i, e := range earnings {
		if !symbols[e.Symbol] {
			continue
		}
		if bySymbol[e.Symbol] == nil {
			bySymbol[e.Symbol] = &pair{}
		}
		p := bySymbol[e.Symbol]

		if e.Date <= today {
			if p.past == nil || e.Date > p.past.Date {
				p.past = &earnings[i]
			}
		} else {
			if p.future == nil || e.Date < p.future.Date {
				p.future = &earnings[i]
			}
		}
	}

	var result []finnhub.Earning
	for _, sym := range sortedKeys(symbols) {
		p := bySymbol[sym]
		if p == nil {
			continue
		}
		if p.past != nil {
			result = append(result, *p.past)
		}
		if p.future != nil {
			result = append(result, *p.future)
		}
	}

	return result
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Simple sort
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}
