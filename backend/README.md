# go-vue-admin 后端服务

Gin + GORM v3 + SQLite + JWT 的 RBAC 后台管理后端，适配前端标准分页与 `{code,msg,data}` 响应结构。

## 启动命令

```bash
cd backend
go mod tidy          # 首次
go run main.go       # 默认 :8080
```

环境变量（可选）：

| 变量 | 默认值 | 说明 |
|---|---|---|
| `PORT` | `8080` | 监听端口 |
| `DB_PATH` | `data/app.db` | SQLite 单文件数据库路径（自动创建目录/文件并迁移） |
| `JWT_SECRET` | `go-vue-admin-jwt-secret` | JWT 签名密钥（生产环境务必覆盖） |

内置超级管理员：`admin / 123456`。首次启动自动执行 `AutoMigrate` 并播种：
super_admin 角色、23 条接口权限、12 条基础菜单，全部绑定给 admin。

## 项目结构

```
backend/
├── main.go                  # 入口：加载配置→初始化DB→注册路由→启动
├── config/config.go         # 环境变量配置
├── database/
│   ├── database.go          # GORM 连接 + AutoMigrate
│   └── seed.go              # 内置 admin/super_admin/菜单/接口权限播种
├── model/model.go           # User/Role/Menu/API + 三张中间表（均软删除）
├── repository/repository.go # User/Role/Menu/API 四个 Repo + 分页
├── service/service.go       # 业务逻辑：登录/JWT/RBAC/个人中心/错误哨兵
├── controller/
│   ├── dto.go               # 全部请求 DTO（binding 校验）
│   ├── auth.go              # 登录/刷新 + 统一错误映射 respond()
│   ├── user.go role.go menu.go api.go profile.go
├── middleware/
│   ├── jwt.go               # Authorization Bearer 校验
│   ├── permission.go        # 接口权限校验（super_admin 直通，白名单，api_permissions 表）
│   ├── cors.go              # 跨域
│   └── recovery.go          # 全局 panic 恢复
├── router/router.go         # 路由注册 + 中间件装配
└── utils/                   # jwt / bcrypt / 统一响应
```

## 响应约定

统一 `{"code":200,"msg":"success","data":...}`；业务错误码：
`200` 成功、`400` 参数/业务校验、`401` 未认证/凭据错误/账号禁用、`403` 无权限、`404` 不存在、`500` 服务器错误。

分页 `data` 结构：`{"total":n,"page":p,"pageSize":s,"list":[...]}`。

## 接口清单

### 认证（公开）

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/login` | 登录 `{username,password}` → `{accessToken,refreshToken,expiresIn,user}` |
| POST | `/api/refresh` | 刷新 `{refreshToken}` → `{accessToken}` |
| GET | `/api/health` | 健康检查 |

### 个人中心（Bearer）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/me` | 当前用户信息（含角色） |
| PUT | `/api/me/password` | 修改密码 `{oldPassword,newPassword}` |
| PUT | `/api/profile` | 更新资料 `{nickname,email,phone,avatar}` |
| GET | `/api/permissions` | 当前用户菜单树/接口权限/权限标识/角色编码 |

### 用户管理（Bearer，接口权限受控）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/users` | 分页，支持 `page,pageSize,keyword,roleId,status` |
| POST | `/api/users` | 创建 `{username,nickname,password,email,phone,avatar,status,roleIds[]}` |
| GET | `/api/users/:id` | 详情 |
| PUT | `/api/users/:id` | 更新 |
| PUT | `/api/users/:id/status` | 启用/禁用 `{status}` |
| DELETE | `/api/users/:id` | 软删除（超管不可删、不可删自己） |

### 角色管理（Bearer，接口权限受控）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/roles` | 分页 `page,pageSize,keyword` |
| GET | `/api/roles/all` | 全量（分配角色下拉） |
| POST | `/api/roles` | 创建 `{name,code,sort,status,remark}` |
| GET | `/api/roles/:id` | 详情（含已分配 menus/apis） |
| PUT | `/api/roles/:id` | 更新（超管角色不可改） |
| DELETE | `/api/roles/:id` | 删除（超管角色不可删，级联解绑） |
| PUT | `/api/roles/:id/menus` | 分配菜单权限 `{menuIds[]}` |
| PUT | `/api/roles/:id/apis` | 分配接口权限 `{apiIds[]}` |

### 菜单管理（Bearer，接口权限受控）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/menus` | 菜单树 |
| GET | `/api/menus/all` | 平铺列表 |
| POST | `/api/menus` | 创建 `{parentId,name,permKey,type(1目录/2菜单/3按钮),icon,url,component,sort,status}` |
| GET | `/api/menus/:id` | 详情 |
| PUT | `/api/menus/:id` | 更新 |
| DELETE | `/api/menus/:id` | 删除（有子节点拒绝） |

### 接口权限（Bearer，接口权限受控）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/apis` | 分页接口权限 |
| GET | `/api/apis/all` | 全量接口权限 |

### 前端联调要点

- 请求头：`Authorization: Bearer <accessToken>`；401 时用 `/api/refresh` 换新 token。
- 路由树/按钮权限用 `/api/permissions` 的 `menus`（树）+ `permKeys`（如 `sys:user:add`）+ `apis`（`path|method`）。
- 前端所有受控请求若 403，提示"无权限"；401 跳登录。
