package router

import (
	"go-vue-admin/controller"
	"go-vue-admin/middleware"
	"go-vue-admin/service"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, cfg service.Deps, db *gorm.DB) {
	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())

	// 健康检查
	r.GET("/api/health", func(c *gin.Context) {
		utils.Success(c, gin.H{"status": "ok"})
	})

	auth := controller.NewAuthController(cfg.Svc)
	user := controller.NewUserController(cfg.Svc)
	role := controller.NewRoleController(cfg.Svc)
	menu := controller.NewMenuController(cfg.Svc)
	api := controller.NewAPIController(cfg.Svc)
	profile := controller.NewProfileController(cfg.Svc)

	// 公开接口
	r.POST("/api/login", auth.Login)
	r.POST("/api/refresh", auth.Refresh)

	// 需要认证的接口
	authed := r.Group("/api",
		middleware.JWT(cfg.JWTSecret),
		middleware.APIPermission(db),
	)
	{
		authed.GET("/me", profile.Me)
		authed.PUT("/me/password", profile.ChangePassword)
		authed.PUT("/profile", profile.UpdateProfile)
		authed.GET("/permissions", profile.Permissions)

		// 用户管理
		authed.GET("/users", user.List)
		authed.GET("/users/:id", user.Get)
		authed.POST("/users", user.Create)
		authed.PUT("/users/:id", user.Update)
		authed.PUT("/users/:id/status", user.SetStatus)
		authed.DELETE("/users/:id", user.Delete)

		// 角色管理
		authed.GET("/roles", role.List)
		authed.GET("/roles/all", role.All)
		authed.GET("/roles/:id", role.Get)
		authed.POST("/roles", role.Create)
		authed.PUT("/roles/:id", role.Update)
		authed.DELETE("/roles/:id", role.Delete)
		authed.PUT("/roles/:id/menus", role.AssignMenus)
		authed.PUT("/roles/:id/apis", role.AssignAPIs)

		// 菜单管理
		authed.GET("/menus", menu.Tree)
		authed.GET("/menus/all", menu.All)
		authed.GET("/menus/:id", menu.Get)
		authed.POST("/menus", menu.Create)
		authed.PUT("/menus/:id", menu.Update)
		authed.DELETE("/menus/:id", menu.Delete)

		// 接口权限
		authed.GET("/apis", api.List)
		authed.GET("/apis/all", api.All)
	}
}
