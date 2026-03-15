package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/oinkywaddles/earnings-cal-cli/internal/finnhub"
)

// PrintListTable prints earnings in Markdown table format.
func PrintListTable(w io.Writer, earnings []finnhub.Earning, from, to, filterDesc string) {
	fmt.Fprintf(w, "Earnings Calendar: %s ~ %s\n", from, to)
	fmt.Fprintf(w, "Source: %s | %d results\n\n", filterDesc, len(earnings))

	if len(earnings) == 0 {
		fmt.Fprintln(w, "No earnings found for this period.")
		return
	}

	// Header
	fmt.Fprintln(w, "| Date       | Symbol | Hour | Quarter   | EPS Est | EPS Act | Rev Est | Rev Act |")
	fmt.Fprintln(w, "|------------|--------|------|-----------|---------|---------|---------|---------|")

	for _, e := range earnings {
		fmt.Fprintf(w, "| %s | %-6s | %-4s | %-9s | %7s | %7s | %7s | %7s |\n",
			e.Date,
			e.Symbol,
			e.Hour,
			FormatQuarter(e.Quarter, e.Year),
			FormatEPS(e.EPSEstimate),
			FormatEPS(e.EPSActual),
			FormatMoney(e.RevenueEstimate),
			FormatMoney(e.RevenueActual),
		)
	}
}

// PrintDetailText prints detail view for one or more earnings.
func PrintDetailText(w io.Writer, earnings []finnhub.Earning, symbolNames map[string]string) {
	for i, e := range earnings {
		if i > 0 {
			fmt.Fprintln(w)
		}

		name := symbolNames[e.Symbol]
		if name == "" {
			name = e.Symbol
		}

		fmt.Fprintf(w, "%s — %s — %s — %s (%s)\n",
			e.Symbol, name, FormatQuarter(e.Quarter, e.Year), e.Date, e.Hour)

		// EPS line
		if e.EPSActual != nil {
			change := FormatChange(e.EPSEstimate, e.EPSActual)
			fmt.Fprintf(w, "  EPS:     Est %s → Actual %s %s\n",
				FormatEPS(e.EPSEstimate), FormatEPS(e.EPSActual), change)
		} else {
			fmt.Fprintf(w, "  EPS:     Est %s → Upcoming\n", FormatEPS(e.EPSEstimate))
		}

		// Revenue line
		if e.RevenueActual != nil {
			change := FormatChange(e.RevenueEstimate, e.RevenueActual)
			fmt.Fprintf(w, "  Revenue: Est %s → Actual %s %s\n",
				FormatMoney(e.RevenueEstimate), FormatMoney(e.RevenueActual), change)
		} else {
			fmt.Fprintf(w, "  Revenue: Est %s → Upcoming\n", FormatMoney(e.RevenueEstimate))
		}
	}
}

// FilterDesc builds the filter description string.
func FilterDesc(symbols []string, all bool, watchlist string) string {
	switch {
	case len(symbols) > 0:
		return "Symbols: " + strings.Join(symbols, ", ")
	case all:
		return "All stocks (no filter)"
	case watchlist != "":
		return "Watchlist: " + watchlist
	default:
		return "S&P 100 ∪ NASDAQ 100 ∪ Dow Jones"
	}
}
