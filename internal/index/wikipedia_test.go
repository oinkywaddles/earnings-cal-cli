package index

import "testing"

func TestIsValidSymbol(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"AAPL", true},
		{"MSFT", true},
		{"BRK.B", true},
		{"BF-B", true},
		{"", false},
		{"TOOLONG", false},
		{"aapl", false},
		{"123", false},
		{"A B", false},
	}
	for _, tt := range tests {
		if got := isValidSymbol(tt.input); got != tt.want {
			t.Errorf("isValidSymbol(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMergeSymbols(t *testing.T) {
	result := mergeSymbols(
		[]string{"AAPL", "MSFT", "GOOGL"},
		[]string{"AAPL", "NVDA"},
		[]string{"MSFT", "JPM"},
	)
	if len(result) != 5 {
		t.Errorf("merged length = %d, want 5", len(result))
	}

	seen := make(map[string]int)
	for _, s := range result {
		seen[s]++
	}
	for s, count := range seen {
		if count != 1 {
			t.Errorf("symbol %s appears %d times, want 1", s, count)
		}
	}
}
