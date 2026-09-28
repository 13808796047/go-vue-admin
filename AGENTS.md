# AGENTS.md（项目根）

仓库根级协作入口。本仓库分前后端：

- `frontend/`：Vue 3 + Vite + Element Plus 后台管理前端（Mock 数据，不接后端）。
  协作文档与自检清单见 **`frontend/AGENTS.md`**，索引：
  - 《设计规范》→ `frontend/docs/design-tokens.md`
  - 《开发规则》→ `frontend/docs/dev-rules.md`
- `backend/`：Go 服务（gva）。本次任务不涉及；改后端前先读 `backend/README.md`。

## 规矩（每次写码前先读，写完对照自检，违反就改）

1. 技术选型锁定：前端只用 Vue 3 + Vite + Element Plus 一个 UI 库，禁止引入第二个。
2. 视觉数值（颜色/圆角/阴影/间距/字号）一律引用 `frontend/src/styles/tokens.scss`，
   主色 `#6366f1`，圆角 12px 体系，间距 8px 体系；禁止魔法数字。
3. 目录结构：`pages / components / composables / styles`（详见《开发规则》目录树）；
   同一 UI 形式出现 ≥2 次必须抽成可复用组件，禁止复制粘贴。
4. 接口契约：组件禁止裸 `fetch`/`axios`，统一走 `src/utils/http.js` + `src/api/*`；
   Mock 兜底可一键关闭（`VITE_ENABLE_MOCK=false`）。
5. 视觉校准：改动视觉后必须跑 Playwright 截图
   （`node scripts/run-visual-check.cjs`，截图输出到 `frontend/screenshots/`），
   对照《设计规范》逐项核对间距/颜色/字号/对齐，列出不一致并逐条修，直到还原度 ≥ 95%。
