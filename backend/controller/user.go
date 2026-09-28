package controller

import (
	"go-vue-admin/middleware"
	"go-vue-admin/repository"
	"go-vue-admin/service"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"strconv"
)

// UserController
type UserController struct {
	Svc *service.Service
}

func NewUserController(svc *service.Service) *UserController {
	return &UserController{Svc: svc}
}

// 从 context 取当前登录用户 ID
func currentUserID(c *gin.Context) uint {
	uidRaw, _ := c.Get(middleware.ContextKeyUserID)
	if uid, ok := uidRaw.(uint); ok {
		return uid
	}
	return 0
}

// 解析路径 :id
func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.Error(c, service.CodeInvalid, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

// List GET /api/users 分页用户列表
func (h *UserController) List(c *gin.Context) {
	var q UserQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	p := &repository.Paginate{Page: q.Page, PageSize: q.PageSize}
	p.Normalize()
	list, total, err := h.Svc.PageUsers(p, q.Keyword, q.RoleID, q.Status)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, utils.PageResult{
		Total:    total,
		Page:     p.Page,
		PageSize: p.PageSize,
		List:     list,
	})
}

// Get GET /api/users/:id
func (h *UserController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	u, err := h.Svc.GetUser(id)
	if err != nil {
		respond(c, err)
		return
	}
	u.Password = ""
	utils.Success(c, u)
}

// Create POST /api/users
func (h *UserController) Create(c *gin.Context) {
	var req UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	if req.Status == 0 {
		req.Status = 1
	}
	u, err := h.Svc.CreateUser(req.Username, req.Nickname, req.Password, req.Email, req.Phone, req.Avatar, req.Status, req.RoleIDs)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, u)
}

// Update PUT /api/users/:id
func (h *UserController) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	u, err := h.Svc.UpdateUser(id, req.Nickname, req.Email, req.Phone, req.Avatar, req.Status, req.RoleIDs)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, u)
}

// SetStatus PUT /api/users/:id/status
func (h *UserController) SetStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UserStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	if err := h.Svc.SetUserStatus(id, req.Status); err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, nil)
}

// Delete DELETE /api/users/:id 软删除
func (h *UserController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if id == currentUserID(c) {
		utils.Error(c, service.CodeInvalid, "不能删除当前登录用户")
		return
	}
	if err := h.Svc.DeleteUser(id); err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, nil)
}
