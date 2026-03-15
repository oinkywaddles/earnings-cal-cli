---
name: earnings-cal
description: Query earnings calendar and key metrics (EPS, revenue) for US stocks
tools: [Bash]
---

# Earnings Calendar CLI

Query earnings calendar data (EPS, revenue estimates/actuals) for S&P 100, NASDAQ 100, and Dow Jones stocks via Finnhub API.

- Setup/install issues → read `skills/SETUP.md`
- Errors during use → read `skills/TROUBLESHOOTING.md`

## Commands

### list — Earnings calendar for a date range

Default shows this week, filtered to major index stocks.

```bash
# Time shortcuts (mutually exclusive, default: this-week)
earnings-cal-cli list                    # this week
earnings-cal-cli list --today
earnings-cal-cli list --tomorrow
earnings-cal-cli list --next-week
earnings-cal-cli list --this-month
earnings-cal-cli list --next-month

# Custom range (max 90 days)
earnings-cal-cli list --from 2026-03-01 --to 2026-03-15

# Filtering
earnings-cal-cli list --symbols AAPL,MSFT   # specific symbols
earnings-cal-cli list --all                  # all stocks, no index filter
earnings-cal-cli list --hour bmo             # bmo (before market open) / amc (after market close) / dmh (during market hours)

# Output control
earnings-cal-cli list --json                 # JSON output
earnings-cal-cli list --limit 20             # cap results
```

**API cost**: 1 Finnhub call per invocation.

### detail — Earnings detail for specific symbols

Returns the most recent past + next upcoming earnings for each symbol.

```bash
earnings-cal-cli detail AAPL
earnings-cal-cli detail AAPL MSFT GOOGL
earnings-cal-cli detail AAPL --from 2025-01-01 --to 2026-03-15
earnings-cal-cli detail AAPL --json
```

**API cost**: 1 Finnhub call per symbol (max 20 symbols).

### init — Pre-cache index constituents

```bash
earnings-cal-cli init              # fetch and cache (skips if fresh)
earnings-cal-cli init --no-cache   # force refresh
```

**API cost**: 0 Finnhub calls. Fetches from Wikipedia only.

## Output formats

Default output is Markdown table (list) or structured text (detail). Use `--json` for machine-readable output.

### list --json schema

```json
{
  "range": { "from": "2026-03-09", "to": "2026-03-15" },
  "filter": "S&P 100 ∪ NASDAQ 100 ∪ Dow Jones",
  "count": 12,
  "earnings": [
    {
      "date": "2026-03-10",
      "symbol": "ORCL",
      "hour": "amc",
      "quarter": 3,
      "year": 2026,
      "epsEstimate": 1.47,
      "epsActual": 1.63,
      "revenueEstimate": 14400000000,
      "revenueActual": 14900000000
    }
  ]
}
```

Fields `epsActual` and `revenueActual` are `null` when earnings have not been reported yet.

### Error JSON

When `--json` is used and an error occurs, output is:

```json
{"error": "error message here"}
```

## Filter priority

```
--symbols flag > EARNINGS_WATCHLIST env > index constituents (S&P 100 ∪ NASDAQ 100 ∪ Dow Jones)
--all skips all filtering
```

## Usage guidelines

- Default table output is human-readable, can be shown to users directly
- Use `--json` only when you need to compute on the data (e.g., filter by EPS surprise)
- Use `--symbols` when the user asks about specific stocks — avoids loading index constituents
- Use `--today` or `--tomorrow` for narrow queries to minimize response size
- Do not pass more than 20 symbols to `detail`
- Custom date ranges are capped at 90 days
