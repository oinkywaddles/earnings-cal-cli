# earnings-cal-cli

CLI tool to query upcoming and past earnings reports for US stocks. Data sourced from [Finnhub](https://finnhub.io/), filtered by S&P 100, NASDAQ 100, and Dow Jones constituents.

## Install

```bash
go install github.com/oinkywaddles/earnings-cal-cli@latest
```

Or download a pre-built binary from [Releases](https://github.com/oinkywaddles/earnings-cal-cli/releases).

## Setup

```bash
# Required: Finnhub API key (free at https://finnhub.io/register)
export FINNHUB_API_KEY=your_key_here

# Optional: pre-cache index constituents (cached 90 days)
earnings-cal-cli init
```

## Usage

### List earnings

```bash
earnings-cal-cli list                          # this week (default)
earnings-cal-cli list --today                  # today only
earnings-cal-cli list --next-week              # next week
earnings-cal-cli list --from 2026-03-01 --to 2026-03-15
earnings-cal-cli list --symbols AAPL,MSFT      # specific symbols
earnings-cal-cli list --all --limit 20         # all stocks, cap 20
earnings-cal-cli list --json                   # JSON output
```

Output:

```
Earnings Calendar: 2026-03-09 ~ 2026-03-15
Source: S&P 100 ∪ NASDAQ 100 ∪ Dow Jones | 3 results

| Date       | Symbol | Hour | Quarter   | EPS Est | EPS Act | Rev Est | Rev Act |
|------------|--------|------|-----------|---------|---------|---------|---------|
| 2026-03-10 | ORCL   | amc  | Q3 FY2026 |    1.47 |    1.63 |  14.4B  |  14.9B  |
| 2026-03-11 | ADBE   | amc  | Q1 FY2026 |    4.97 |       — |  5.63B  |       — |
```

### Detail view

```bash
earnings-cal-cli detail AAPL                   # last + next earnings
earnings-cal-cli detail AAPL MSFT GOOGL        # multiple symbols
earnings-cal-cli detail AAPL --json
```

Output:

```
ORCL — ORCL — Q3 FY2026 — 2026-03-10 (amc)
  EPS:     Est 1.47 → Actual 1.63 (+10.9%)
  Revenue: Est 14.4B → Actual 14.9B (+3.5%)
```

### Filtering

```
--symbols flag > EARNINGS_WATCHLIST env > index constituents
--all skips all filtering
--hour bmo|amc|dmh filters by reporting time
```

## Environment variables

| Variable | Required | Description |
|----------|----------|-------------|
| `FINNHUB_API_KEY` | Yes | API key from [Finnhub](https://finnhub.io/register) |
| `EARNINGS_WATCHLIST` | No | Default symbol filter, e.g. `AAPL,MSFT,GOOGL` |

## License

[MIT](LICENSE)
