package controller

import (
	"go-vue-admin/service"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

// ProfileController 个人中心 + 权限
type ProfileController struct {
	Svc *service.Service
}

func NewProfileController(svc *service.Service) *ProfileController {
	return &ProfileController{Svc: svc}
}

// Me GET /api/me 当前用户信息（含角色）
func (h *ProfileController) Me(c *gin.Context) {
	uid := currentUserID(c)
	u, err := h.Svc.GetProfile(uid)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, u)
}

// UpdateProfile PUT /api/profile
func (h *ProfileController) UpdateProfile(c *gin.Context) {
	uid := currentUserID(c)
	var req ProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	u, err := h.Svc.UpdateProfile(uid, req.Nickname, req.Email, req.Phone, req.Avatar)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, u)
}

// ChangePassword PUT /api/me/password
func (h *ProfileController) ChangePassword(c *gin.Context) {
	uid := currentUserID(c)
	var req PasswordChangeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	if err := h.Svc.ChangePassword(uid, req.OldPassword, req.NewPassword); err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, nil)
}

// Permissions GET /api/permissions 当前用户菜单/接口权限
func (h *ProfileController) Permissions(c *gin.Context) {
	uid := currentUserID(c)
	perms, err := h.Svc.UserPermissions(uid)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, perms)
}
