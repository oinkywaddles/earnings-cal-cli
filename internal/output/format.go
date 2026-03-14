package output

import (
	"fmt"
	"math"
)

// FormatMoney formats a number with B/M/K suffix.
func FormatMoney(v *float64) string {
	if v == nil {
		return "—"
	}
	n := *v
	abs := math.Abs(n)

	switch {
	case abs >= 1e12:
		return fmt.Sprintf("%.1fT", n/1e12)
	case abs >= 1e9:
		return fmt.Sprintf("%.1fB", n/1e9)
	case abs >= 1e6:
		return fmt.Sprintf("%.1fM", n/1e6)
	case abs >= 1e3:
		return fmt.Sprintf("%.1fK", n/1e3)
	default:
		return fmt.Sprintf("%.2f", n)
	}
}

// FormatEPS formats an EPS value.
func FormatEPS(v *float64) string {
	if v == nil {
		return "—"
	}
	return fmt.Sprintf("%.2f", *v)
}

// FormatQuarter returns something like "Q3 FY2026".
func FormatQuarter(quarter, year int) string {
	return fmt.Sprintf("Q%d FY%d", quarter, year)
}

// FormatChange calculates percentage change from estimate to actual.
func FormatChange(est, act *float64) string {
	if est == nil || act == nil || *est == 0 {
		return ""
	}
	pct := (*act - *est) / math.Abs(*est) * 100
	sign := "+"
	if pct < 0 {
		sign = ""
	}
	return fmt.Sprintf("(%s%.1f%%)", sign, pct)
}
