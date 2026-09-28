package middleware

import (
	"net/http"

	"go-vue-admin/model"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIPermission 接口权限中间件：校验当前用户的角色是否授权 请求 method+path
// - 超级管理员(super_admin)直接放行
// - 登录/刷新/自身信息白名单放行
// - 其余按 api_permissions 表（经 role_apis -> user_roles 关联）校验
func APIPermission(db *gorm.DB) gin.HandlerFunc {
	// 白名单：自身信息与认证接口
	whitelist := map[string]bool{
		"GET|/api/me":          true,
		"PUT|/api/me":          true,
		"PUT|/api/me/password": true,
		"PUT|/api/profile":     true,
		"GET|/api/permissions": true,
	}

	return func(c *gin.Context) {
		uidRaw, exists := c.Get(ContextKeyUserID)
		if !exists {
			utils.Error(c, http.StatusUnauthorized, "未认证")
			c.Abort()
			return
		}
		uid := uidRaw.(uint)

		var u model.User
		if err := db.Preload("Roles").First(&u, uid).Error; err != nil {
			utils.Error(c, http.StatusUnauthorized, "用户不存在")
			c.Abort()
			return
		}

		// 超级管理员放行
		for _, ro := range u.Roles {
			if ro.Code == "super_admin" {
				c.Next()
				return
			}
		}

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		key := c.Request.Method + "|" + path

		// 白名单
		if whitelist[key] {
			c.Next()
			return
		}

		// 该用户是否被授权此 method+path
		var count int64
		err := db.Model(&model.API{}).
			Joins("JOIN role_apis ON role_apis.api_id = api_permissions.id").
			Joins("JOIN user_roles ON user_roles.role_id = role_apis.role_id").
			Where("user_roles.user_id = ? AND api_permissions.path = ? AND api_permissions.method = ?",
				uid, path, c.Request.Method).
			Count(&count).Error
		if err == nil && count > 0 {
			c.Next()
			return
		}

		utils.Error(c, http.StatusForbidden, "无权限访问该接口: "+key)
		c.Abort()
	}
}
