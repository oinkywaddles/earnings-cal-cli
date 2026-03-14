package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/oinkywaddles/earnings-cal-cli/internal/finnhub"
)

func TestPrintListJSON(t *testing.T) {
	earnings := []finnhub.Earning{
		{
			Date:        "2026-03-10",
			Symbol:      "ORCL",
			Hour:        "amc",
			Quarter:     3,
			Year:        2026,
			EPSEstimate: f64(1.47),
			EPSActual:   f64(1.63),
		},
	}
	var buf bytes.Buffer
	PrintListJSON(&buf, earnings, "2026-03-09", "2026-03-15", "test")

	var result jsonListOutput
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("count = %d, want 1", result.Count)
	}
	if result.Range.From != "2026-03-09" {
		t.Errorf("range.from = %q, want 2026-03-09", result.Range.From)
	}
	if result.Earnings[0].Symbol != "ORCL" {
		t.Errorf("symbol = %q, want ORCL", result.Earnings[0].Symbol)
	}
}

func TestPrintDetailJSON(t *testing.T) {
	earnings := []finnhub.Earning{
		{Date: "2026-03-10", Symbol: "AAPL"},
	}
	var buf bytes.Buffer
	PrintDetailJSON(&buf, earnings)

	var result []finnhub.Earning
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result) != 1 || result[0].Symbol != "AAPL" {
		t.Errorf("unexpected result: %+v", result)
	}
}
