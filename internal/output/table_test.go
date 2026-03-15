package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/oinkywaddles/earnings-cal-cli/internal/finnhub"
)

func TestPrintListTable_Empty(t *testing.T) {
	var buf bytes.Buffer
	PrintListTable(&buf, nil, "2026-03-09", "2026-03-15", "test")
	out := buf.String()
	if !strings.Contains(out, "No earnings found") {
		t.Errorf("expected 'No earnings found' in output, got: %s", out)
	}
	if !strings.Contains(out, "0 results") {
		t.Errorf("expected '0 results' in output, got: %s", out)
	}
}

func TestPrintListTable_WithData(t *testing.T) {
	earnings := []finnhub.Earning{
		{
			Date:            "2026-03-10",
			Symbol:          "ORCL",
			Hour:            "amc",
			Quarter:         3,
			Year:            2026,
			EPSEstimate:     f64(1.47),
			EPSActual:       f64(1.63),
			RevenueEstimate: f64(14400000000),
			RevenueActual:   f64(14900000000),
		},
	}
	var buf bytes.Buffer
	PrintListTable(&buf, earnings, "2026-03-09", "2026-03-15", "test")
	out := buf.String()

	if !strings.Contains(out, "1 results") {
		t.Errorf("expected '1 results', got: %s", out)
	}
	if !strings.Contains(out, "ORCL") {
		t.Errorf("expected 'ORCL' in table, got: %s", out)
	}
	if !strings.Contains(out, "1.63") {
		t.Errorf("expected EPS actual '1.63' in table, got: %s", out)
	}
}

func TestPrintDetailText_Upcoming(t *testing.T) {
	earnings := []finnhub.Earning{
		{
			Date:            "2026-03-11",
			Symbol:          "ADBE",
			Hour:            "amc",
			Quarter:         1,
			Year:            2026,
			EPSEstimate:     f64(4.97),
			RevenueEstimate: f64(5630000000),
		},
	}
	var buf bytes.Buffer
	PrintDetailText(&buf, earnings, nil)
	out := buf.String()

	if !strings.Contains(out, "Upcoming") {
		t.Errorf("expected 'Upcoming' for unreported earnings, got: %s", out)
	}
	if !strings.Contains(out, "4.97") {
		t.Errorf("expected EPS estimate '4.97', got: %s", out)
	}
}

func TestFilterDesc(t *testing.T) {
	if got := FilterDesc([]string{"AAPL"}, false, ""); got != "Symbols: AAPL" {
		t.Errorf("got %q", got)
	}
	if got := FilterDesc(nil, true, ""); got != "All stocks (no filter)" {
		t.Errorf("got %q", got)
	}
	if got := FilterDesc(nil, false, "AAPL,MSFT"); got != "Watchlist: AAPL,MSFT" {
		t.Errorf("got %q", got)
	}
	if got := FilterDesc(nil, false, ""); !strings.Contains(got, "S&P 100") {
		t.Errorf("got %q", got)
	}
}
