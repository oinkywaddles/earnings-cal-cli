# TODO

## P2 — Extended Data

- [ ] `/stock/financials-reported` — SEC 申报的完整财务数据（利润率、现金流等）
- [ ] `/stock/earnings` — 历史季度 EPS surprise 数据
- [ ] `/company-news` — 公司新闻（含财报相关）
- [ ] 管理层指引 (guidance) 数据

## P2 — Features

- [ ] `watch` 命令 — 监听即将发布的财报并提醒
- [ ] 自定义 watchlist 文件（替代环境变量）
- [ ] 输出排序选项（按 symbol / EPS surprise / revenue）
- [ ] `surprise` 子命令 — 显示 beat/miss 排行

## P2 — Infrastructure

- [ ] 增加 cmd/ 包的集成测试
- [ ] 增加 finnhub client 的 mock 测试
- [ ] Homebrew tap 发布
