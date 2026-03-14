package output

import (
	"testing"
)

func f64(v float64) *float64 { return &v }

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		input *float64
		want  string
	}{
		{nil, "—"},
		{f64(0), "0.00"},
		{f64(500), "500.00"},
		{f64(1500), "1.5K"},
		{f64(14400000000), "14.4B"},
		{f64(5630000000), "5.6B"},
		{f64(2500000), "2.5M"},
		{f64(1200000000000), "1.2T"},
		{f64(-5630000000), "-5.6B"},
	}
	for _, tt := range tests {
		got := FormatMoney(tt.input)
		if got != tt.want {
			t.Errorf("FormatMoney(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatEPS(t *testing.T) {
	tests := []struct {
		input *float64
		want  string
	}{
		{nil, "—"},
		{f64(1.47), "1.47"},
		{f64(-0.5), "-0.50"},
		{f64(0), "0.00"},
	}
	for _, tt := range tests {
		got := FormatEPS(tt.input)
		if got != tt.want {
			t.Errorf("FormatEPS(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatQuarter(t *testing.T) {
	if got := FormatQuarter(3, 2026); got != "Q3 FY2026" {
		t.Errorf("FormatQuarter(3, 2026) = %q, want %q", got, "Q3 FY2026")
	}
}

func TestFormatChange(t *testing.T) {
	tests := []struct {
		est, act *float64
		want     string
	}{
		{nil, f64(1.0), ""},
		{f64(1.0), nil, ""},
		{f64(0), f64(1.0), ""},
		{f64(1.47), f64(1.63), "(+10.9%)"},
		{f64(10.0), f64(8.0), "(-20.0%)"},
		{f64(1.0), f64(1.0), "(+0.0%)"},
	}
	for _, tt := range tests {
		got := FormatChange(tt.est, tt.act)
		if got != tt.want {
			t.Errorf("FormatChange(%v, %v) = %q, want %q", tt.est, tt.act, got, tt.want)
		}
	}
}
