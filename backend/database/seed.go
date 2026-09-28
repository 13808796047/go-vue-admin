package database

import (
	"go-vue-admin/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seed 内置超级管理员 admin/123456 及基础 RBAC 数据
func seed(db *gorm.DB) {
	var uCount int64
	db.Model(&model.User{}).Count(&uCount)
	if uCount > 0 {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), 10)
	admin := model.User{Username: "admin", Nickname: "超级管理员", Email: "admin@example.com", Password: string(hash), Status: 1}
	db.Create(&admin)

	superRole := model.Role{Name: "超级管理员", Code: "super_admin", Sort: 0, Status: 1, Remark: "内置角色"}
	db.Create(&superRole)

	// 内置 API 权限
	apiDefs := []struct {
		Path, Method, Name, Desc string
	}{
		{"/api/login", "POST", "登录", "用户登录"},
		{"/api/refresh", "POST", "刷新令牌", "JWT刷新"},
		{"/api/me", "GET", "个人信息", "获取当前登录用户"},
		{"/api/me/password", "PUT", "修改密码", "个人中心修改密码"},
		{"/api/profile", "PUT", "更新资料", "个人中心更新资料"},
		{"/api/users", "GET", "用户列表", "分页查询用户"},
		{"/api/users", "POST", "创建用户", "新增用户"},
		{"/api/users/:id", "GET", "用户详情", "查看用户"},
		{"/api/users/:id", "PUT", "更新用户", "编辑用户"},
		{"/api/users/:id", "DELETE", "删除用户", "软删除用户"},
		{"/api/roles", "GET", "角色列表", "分页查询角色"},
		{"/api/roles", "POST", "创建角色", "新增角色"},
		{"/api/roles/:id", "GET", "角色详情", "查看角色"},
		{"/api/roles/:id", "PUT", "更新角色", "编辑角色"},
		{"/api/roles/:id", "DELETE", "删除角色", "删除角色"},
		{"/api/roles/:id/menus", "PUT", "分配菜单", "角色分配菜单权限"},
		{"/api/roles/:id/apis", "PUT", "分配接口", "角色分配接口权限"},
		{"/api/menus", "GET", "菜单树", "获取菜单树"},
		{"/api/menus", "POST", "创建菜单", "新增菜单/权限"},
		{"/api/menus/:id", "GET", "菜单详情", "查看菜单"},
		{"/api/menus/:id", "PUT", "更新菜单", "编辑菜单"},
		{"/api/menus/:id", "DELETE", "删除菜单", "删除菜单"},
		{"/api/permissions", "GET", "用户权限", "获取当前用户菜单/接口权限"},
	}
	apis := make([]model.API, 0, len(apiDefs))
	for _, a := range apiDefs {
		apis = append(apis, model.API{Path: a.Path, Method: a.Method, Name: a.Name, Desc: a.Desc})
	}
	db.Create(&apis)

	// 内置菜单
	menuDefs := []model.Menu{
		{ParentID: 0, Name: "系统管理", PermKey: "sys", Type: 1, Icon: "setting", URL: "/system", Component: "SystemLayout", Sort: 1},
		{ParentID: 1, Name: "用户管理", PermKey: "sys:user", Type: 2, Icon: "user", URL: "/system/user", Component: "system/User", Sort: 1},
		{ParentID: 1, Name: "角色管理", PermKey: "sys:role", Type: 2, Icon: "team", URL: "/system/role", Component: "system/Role", Sort: 2},
		{ParentID: 1, Name: "菜单管理", PermKey: "sys:menu", Type: 2, Icon: "apartment", URL: "/system/menu", Component: "system/Menu", Sort: 3},
		{ParentID: 2, Name: "用户新增", PermKey: "sys:user:add", Type: 3, Sort: 1},
		{ParentID: 2, Name: "用户编辑", PermKey: "sys:user:edit", Type: 3, Sort: 2},
		{ParentID: 2, Name: "用户删除", PermKey: "sys:user:delete", Type: 3, Sort: 3},
		{ParentID: 3, Name: "角色新增", PermKey: "sys:role:add", Type: 3, Sort: 1},
		{ParentID: 3, Name: "角色删除", PermKey: "sys:role:delete", Type: 3, Sort: 2},
		{ParentID: 4, Name: "菜单新增", PermKey: "sys:menu:add", Type: 3, Sort: 1},
		{ParentID: 4, Name: "菜单删除", PermKey: "sys:menu:delete", Type: 3, Sort: 2},
		{ParentID: 0, Name: "个人中心", PermKey: "profile", Type: 2, Icon: "idcard", URL: "/profile", Component: "profile/Index", Sort: 2},
	}
	db.Create(&menuDefs)

	// 超级管理员：拥有全部菜单与接口权限，并绑定 admin 用户
	menuRows := make([]model.RoleMenus, 0, len(menuDefs))
	for _, m := range menuDefs {
		menuRows = append(menuRows, model.RoleMenus{RoleID: superRole.ID, MenuID: m.ID})
	}
	apiRows := make([]model.RoleApis, 0, len(apis))
	for _, a := range apis {
		apiRows = append(apiRows, model.RoleApis{RoleID: superRole.ID, ApiID: a.ID})
	}
	db.Create(&menuRows)
	db.Create(&apiRows)
	db.Create(&model.UserRoles{UserID: admin.ID, RoleID: superRole.ID})
}
