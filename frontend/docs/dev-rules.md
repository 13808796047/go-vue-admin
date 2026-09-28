# 《开发规则》

本前端项目的结构与协作约定。新增/修改代码前对照，写完自检。

## 目录结构

```
frontend/
├─ AGENTS.md              # 协作入口 + 自检清单（本文档索引）
├─ docs/
│  ├─ design-tokens.md    # 《设计规范》
│  └─ dev-rules.md        # 《开发规则》（本文档）
├─ index.html
├─ vite.config.js
└─ src/
   ├─ main.js             # 入口：挂载 App、注册 Element Plus / 路由 / 全局样式
   ├─ App.vue            # 根组件
   ├─ api/               # 接口层：每个后端资源一个文件，只导 http 调用
   │  └─ auth.js, user.js, role.js, menu.js, permission.js
   ├─ components/        # 可复用组件（见下"复用阈值"）
   │  ├─ layout/          # 布局：AdminLayout、Sidebar、Header
   │  ├─ base/            # 原子：PageHeader、PageCard、StatusTag、EmptyState、RowActions、CodeChip、RoleTags
   │  └─ data/            # 数据交互：Toolbar、FormDrawer
   ├─ composables/       # 跨组件逻辑：useAuth、useTable、usePermission
   ├─ pages/             # 路由页面（一个功能一个文件）
   │  ├─ login/LoginPage.vue
   │  ├─ user/UserManagePage.vue
   │  ├─ role/RoleManagePage.vue
   │  ├─ menu/MenuManagePage.vue
   │  ├─ permission/PermissionPage.vue
   │  └─ profile/ProfilePage.vue
   ├─ router/            # vue-router 配置 + 登录态/权限守卫（index.js）
   ├─ styles/            # 全局样式 + tokens.scss（视觉基准唯一来源）
   │  ├─ tokens.scss
   │  ├─ reset.scss
   │  └─ base.scss
   ├─ mock/              # mock 拦截器（开发兜底，可用开关关闭）
   │  └─ index.js
   ├─ utils/
   │  ├─ http.js         # axios 实例 + 请求/响应拦截 + 统一错误
   │  └─ storage.js     # token / 本地存储封装
   └─ constants/        # 枚举、字典（如角色类型、状态码）
```

## 复用阈值（强制）

- **同一 UI 形式出现 ≥2 次 → 抽成可复用组件**，禁止复制粘贴。
  判断单位是"视觉单元"：一张带标题+操作区的卡片、一行带状态标签的表格行、一个搜索+筛选的表单头，重复即抽。
- 组件归属：
  - 跨页通用（按钮/输入/标签/空态/卡片壳）→ `components/base/`
  - 布局骨架 → `components/layout/`
  - 表格/抽屉/表单等数据交互 → `components/data/`
  - 纯逻辑（请求封装、权限判断、表格通用逻辑）→ `composables/`，前缀 `use`
- 一个组件若只在一个页面用且无复用价值，可留在该 page 内；但**只要结构在第二处出现就外提**。

## 接口契约（对接后端）

- 所有 HTTP 走 `src/utils/http.js` 的 axios 实例，**禁止在组件里直接 `axios`/`fetch`**。
- 每个资源一个 `src/api/<resource>.js`，只导出语义化函数（`login`、`listUsers`、`createRole`…），不导出裸 URL。
- 请求前缀 `/api`，统一 `{ code, message, data }` 返回体；`code === 0` 为成功。
- **mock 兜底**：`src/mock/` 注册拦截器模拟上述接口，开发态默认开启；
  一键关闭：`VITE_ENABLE_MOCK=false`（写 `.env`）或代码里 `setMockEnabled(false)`。后端就绪后切真实接口不改调用方。
- 鉴权：登录拿到 `token` 存 `localStorage`，`http.js` 响应拦截自动带 `Authorization: Bearer`；
  收到 `401` 跳登录。

## 命名

- 页面文件 / 组件：`PascalCase`（`UserManagePage.vue`、`DataTable.vue`）。
- 路由 path：`kebab-case`（`/user-manage`、`/menu-manage`）。
- composable：`usePascal`（`useAuth`、`usePermission`）。
- 组件内变量 / 事件：`camelCase`；自定义事件 emit `kebab-case`。
- 常量枚举：`UPPER_SNAKE`（`ROLE_TYPES`）。

## 代码风格

- `<script setup>` 组合式 API；组件用 `<script setup>`，不要 Options API。
- 样式 scoped；视觉数值只引用 tokens（见《设计规范》），不在组件里写魔法值。
- Element Plus 全量注册（`main.js` 中 `app.use(ElementPlus)` + 中文 locale），模板里用 `el-*` 组件；
  图标用 `@element-plus/icons-vue` 按需 import。全项目仅此一个 UI 库。
- 一个文件一个功能；函数职责单一；无未使用导入。
- 文本界面用简体中文；标识符、命令、API 名保持英文。

## 权限与路由守卫

- `router/index.js`：未登录跳 `/login`；刷新后通过 `useAuth.ensurePermissions()` 恢复权限码；
  路由 `meta.perm` 存在但未持有时重定向到 `/profile`。侧边栏由 `usePermission().visibleMenus` 过滤。

## 自检（与 AGENTS.md 同步）

写完任一文件，回 AGENTS.md 的 Self-check 逐条过；涉及结构/复用/接口三条重点核查。
