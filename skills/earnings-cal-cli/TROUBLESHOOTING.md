# Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `FINNHUB_API_KEY not set` | Environment variable missing | `export FINNHUB_API_KEY=...` (get key at https://finnhub.io/register) |
| `finnhub rate limit exceeded` | Over 60 calls/min (free tier) | Wait 1 minute, reduce call frequency |
| `finnhub returned status 403` | Invalid API key | Verify key at https://finnhub.io/dashboard |
| `Failed to fetch index constituents` | Wikipedia unreachable | Check network, or use `--symbols` / `--all` to bypass |
| `date range exceeds 90 days` | `--from`/`--to` span too large | Use a narrower range or time shortcuts (`--this-month`, etc.) |
| `too many symbols (N), max is 20` | Too many args to `detail` | Split into multiple `detail` calls |
| `no earnings found` | Symbol has no earnings in range | Verify symbol is correct, try wider date range |
| `command not found: earnings-cal-cli` | Binary not installed | See `skills/SETUP.md` for installation |
