# Setup Guide

## Install

```bash
go install github.com/oinkywaddles/earnings-cal-cli@latest
```

Or build from source:

```bash
git clone https://github.com/oinkywaddles/earnings-cal-cli.git
cd earnings-cal-cli
make build
# binary: ./earnings-cal-cli
```

## Configuration

### FINNHUB_API_KEY (required)

Get a free API key at https://finnhub.io/register, then set it:

```bash
export FINNHUB_API_KEY=your_key_here
```

Free tier limit: 60 API calls/minute.

### EARNINGS_WATCHLIST (optional)

Set a default symbol filter:

```bash
export EARNINGS_WATCHLIST=AAPL,MSFT,GOOGL,NVDA
```

When set, `list` uses this as the default filter instead of index constituents.

## Initialize index cache

Pre-cache S&P 500, NASDAQ 100, and Dow Jones constituents (fetched from Wikipedia, cached 90 days):

```bash
earnings-cal-cli init
```

This step is optional — `list` will auto-fetch on first use if needed. Running `init` avoids the delay on the first `list` call.

## Verify

```bash
earnings-cal-cli list --today --json
```

If you encounter errors, read `skills/TROUBLESHOOTING.md`.
