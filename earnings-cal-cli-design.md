# earnings-cal-cli 设计方案

## 概述

财报日历 CLI 工具，供 Claude Code skill 调用 + 独立开源项目。

- **语言**: Go（单 binary，零运行时依赖）
- **CLI 框架**: cobra
- **数据源**: Finnhub API（免费 tier）+ Wikipedia（指数成分股）
- **开源协议**: MIT

---

## CLI 接口

### `list` — 财报日历列表

```bash
# 时间快捷方式（互斥，默认 this-week）
earnings-cal-cli list
earnings-cal-cli list --today
earnings-cal-cli list --tomorrow
earnings-cal-cli list --this-week
earnings-cal-cli list --next-week
earnings-cal-cli list --this-month
earnings-cal-cli list --next-month

# 自定义日期范围
earnings-cal-cli list --from 2026-03-01 --to 2026-03-15

# 过滤
earnings-cal-cli list --symbols AAPL,MSFT    # 指定 symbol
earnings-cal-cli list --all                  # 不做指数过滤，全量
earnings-cal-cli list --hour bmo             # bmo / amc / dmh

# 输出与数量
earnings-cal-cli list --json
earnings-cal-cli list --limit 20
```

### `detail` — 单股/批量详情

```bash
earnings-cal-cli detail AAPL
earnings-cal-cli detail AAPL MSFT GOOGL
earnings-cal-cli detail AAPL --from 2025-01-01 --to 2026-03-13
earnings-cal-cli detail AAPL --json
```

### 全局 flags

```
--json          JSON 输出
--no-cache      强制刷新缓存
--help / -h     帮助
--version / -v  版本
```

---

## 输出格式

### `list` 默认（Markdown 表格）

```
Earnings Calendar: 2026-03-09 ~ 2026-03-15 (this week)
Source: S&P 500 ∪ NASDAQ 100 ∪ Dow Jones | 12 results

| Date       | Symbol | Hour | Quarter   | EPS Est | EPS Act | Rev Est | Rev Act |
|------------|--------|------|-----------|---------|---------|---------|---------|
| 2026-03-10 | ORCL   | amc  | Q3 FY2026 |    1.47 |    1.63 |  14.4B  |  14.9B  |
| 2026-03-11 | ADBE   | amc  | Q1 FY2026 |    4.97 |       — |  5.63B  |       — |
```

- 已发布：Est 和 Act 都有值
- 未发布：Act 显示 `—`
- 金额自动格式化（B/M/K）

### `detail` 默认

```
ORCL — Oracle Corp — Q3 FY2026 — 2026-03-10 (amc)
  EPS:     Est 1.47 → Actual 1.63 (+10.9%)
  Revenue: Est 14.4B → Actual 14.9B (+3.5%)

ADBE — Adobe Inc — Q1 FY2026 — 2026-03-11 (amc)
  EPS:     Est 4.97 → Upcoming
  Revenue: Est 5.63B → Upcoming
```

### `--json`

```json
{
  "range": { "from": "2026-03-09", "to": "2026-03-15" },
  "filter": "index",
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

---

## 配置

| 配置项 | 方式 | 必填 | 说明 |
|--------|------|------|------|
| API Key | `FINNHUB_API_KEY` 环境变量 | 是 | 缺失时报错并提示 |
| Watchlist | `EARNINGS_WATCHLIST` 环境变量 | 否 | `AAPL,MSFT,GOOGL` |
| 缓存目录 | `~/.cache/earnings-cal-cli/` | — | 自动创建 |

### 过滤优先级

```
--symbols 参数 > EARNINGS_WATCHLIST 环境变量 > 三大指数并集
--all 跳过所有过滤
```

---

## 指数成分股

### 数据源：Wikipedia

| 指数 | 页面 |
|------|------|
| S&P 500 | `en.wikipedia.org/wiki/List_of_S%26P_500_companies` |
| NASDAQ 100 | `en.wikipedia.org/wiki/Nasdaq-100` |
| Dow Jones | `en.wikipedia.org/wiki/Dow_Jones_Industrial_Average` |

### 缓存

- 文件：`~/.cache/earnings-cal-cli/index_constituents.json`
- 过期：7 天
- `--no-cache` 强制刷新

```json
{
  "updated_at": "2026-03-14T00:00:00Z",
  "sp500": ["AAPL", "MSFT", ...],
  "nasdaq100": ["AAPL", "NVDA", ...],
  "dowjones": ["AAPL", "MSFT", ...],
  "merged": ["AAPL", "MSFT", "NVDA", ...]
}
```

---

## 项目结构

```
earnings-cal-cli/
├── cmd/
│   ├── root.go              # cobra root, 全局 flags
│   ├── list.go              # list subcommand
│   └── detail.go            # detail subcommand
├── internal/
│   ├── finnhub/
│   │   └── client.go        # Finnhub API 封装
│   ├── index/
│   │   ├── wikipedia.go     # Wikipedia HTML 解析
│   │   └── cache.go         # 本地缓存
│   └── output/
│       ├── table.go         # Markdown 表格
│       ├── json.go          # JSON 输出
│       └── format.go        # 金额格式化工具
├── main.go
├── go.mod
├── go.sum
├── Makefile
├── .goreleaser.yaml         # 多平台发布
├── LICENSE
└── README.md
```

---

## Go 依赖

| 依赖 | 用途 |
|------|------|
| `github.com/spf13/cobra` | CLI 框架 |
| `golang.org/x/net/html` | HTML 解析 |
| 标准库 `net/http`, `encoding/json`, `os`, `time` | HTTP/JSON/缓存 |

最小依赖原则，不引入重型库。

---

## 核心流程

### `list`

```
1. 解析时间参数 → from/to
2. 确定过滤：--symbols / --all / EARNINGS_WATCHLIST / 指数成分
3. Finnhub GET /calendar/earnings?from=X&to=Y
4. 本地过滤：symbol + hour
5. 排序：date 升序
6. 输出：table 或 json
```

### `detail`

```
1. 拉日期范围全量（免费 tier 按 symbol 查可能返回空）
2. 本地按 symbol 过滤
3. 格式化详情输出
```

---

## 分发

```bash
# go install
go install github.com/<user>/earnings-cal-cli@latest

# GitHub Release（goreleaser 生成多平台 binary）
# darwin/linux/windows × amd64/arm64

# Homebrew（后续）
brew install <user>/tap/earnings-cal-cli
```

### CI/CD

- push to main → tests + lint
- tag push → goreleaser 发布

---

## 实现优先级

| 阶段 | 内容 |
|------|------|
| **P0** | `list` + Finnhub 调用 + 时间快捷方式 + Markdown 输出 |
| **P0** | Wikipedia 成分股解析 + 缓存 |
| **P1** | `detail` 命令 |
| **P1** | `--json` / `--symbols` / `EARNINGS_WATCHLIST` / `--hour` / `--all` |
| **P2** | goreleaser + CI/CD |
| **P2** | Homebrew tap |
| **P2** | Claude Code SKILL.md |

---

## 验证

```bash
go build -o earnings-cal-cli .

export FINNHUB_API_KEY=xxx
./earnings-cal-cli list
./earnings-cal-cli list --today
./earnings-cal-cli list --json
./earnings-cal-cli list --symbols AAPL,MSFT
./earnings-cal-cli list --all --limit 10
./earnings-cal-cli detail AAPL
./earnings-cal-cli detail AAPL --json

# 缓存
ls ~/.cache/earnings-cal-cli/
./earnings-cal-cli list --no-cache
```
