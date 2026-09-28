package controller

import (
	"errors"

	"go-vue-admin/service"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

// respond 将 service 层错误映射为统一 {code,msg,data} 响应
func respond(c *gin.Context, err error) {
	if err == nil {
		return
	}
	switch {
	case errors.Is(err, service.ErrNotFound):
		utils.Error(c, service.CodeNotFound, err.Error())
	case errors.Is(err, service.ErrAlreadyExists), errors.Is(err, service.ErrHasChildren):
		utils.Error(c, service.CodeInvalid, err.Error())
	case errors.Is(err, service.ErrInvalidCred), errors.Is(err, service.ErrOldPassword),
		errors.Is(err, service.ErrAccountLocked):
		utils.Error(c, service.CodeUnauth, err.Error())
	case errors.Is(err, service.ErrNoPermission), errors.Is(err, service.ErrProtected):
		utils.Error(c, service.CodeNoPerm, err.Error())
	default:
		utils.Error(c, service.CodeServer, err.Error())
	}
}

// AuthController
type AuthController struct {
	Svc *service.Service
}

func NewAuthController(svc *service.Service) *AuthController {
	return &AuthController{Svc: svc}
}

// Login POST /api/login
func (h *AuthController) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	res, err := h.Svc.Login(req.Username, req.Password)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, res)
}

// Refresh POST /api/refresh
func (h *AuthController) Refresh(c *gin.Context) {
	var req RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	token, err := h.Svc.Refresh(req.RefreshToken)
	if err != nil {
		utils.Error(c, service.CodeUnauth, err.Error())
		return
	}
	utils.Success(c, gin.H{"accessToken": token})
}
