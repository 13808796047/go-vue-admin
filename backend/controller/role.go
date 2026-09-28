package controller

import (
	"go-vue-admin/repository"
	"go-vue-admin/service"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

// RoleController
type RoleController struct {
	Svc *service.Service
}

func NewRoleController(svc *service.Service) *RoleController {
	return &RoleController{Svc: svc}
}

// List GET /api/roles
func (h *RoleController) List(c *gin.Context) {
	var q RoleQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	p := &repository.Paginate{Page: q.Page, PageSize: q.PageSize}
	p.Normalize()
	list, total, err := h.Svc.PageRoles(p, q.Keyword)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, utils.PageResult{Total: total, Page: p.Page, PageSize: p.PageSize, List: list})
}

// All GET /api/roles/all 全量（用于用户分配角色下拉）
func (h *RoleController) All(c *gin.Context) {
	p := &repository.Paginate{Page: 1, PageSize: 100}
	p.Normalize()
	list, _, err := h.Svc.PageRoles(p, "")
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, list)
}

// Get GET /api/roles/:id
func (h *RoleController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	ro, err := h.Svc.GetRole(id)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, ro)
}

// Create POST /api/roles
func (h *RoleController) Create(c *gin.Context) {
	var req RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	if req.Status == 0 {
		req.Status = 1
	}
	ro, err := h.Svc.CreateRole(req.Name, req.Code, req.Sort, req.Status, req.Remark)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, ro)
}

// Update PUT /api/roles/:id
func (h *RoleController) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	ro, err := h.Svc.UpdateRole(id, req.Name, req.Sort, req.Status, req.Remark)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, ro)
}

// Delete DELETE /api/roles/:id
func (h *RoleController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.Svc.DeleteRole(id); err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, nil)
}

// AssignMenus PUT /api/roles/:id/menus 角色分配菜单
func (h *RoleController) AssignMenus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req RoleMenuAssignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	if err := h.Svc.AssignRoleMenus(id, req.MenuIDs); err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, nil)
}

// AssignAPIs PUT /api/roles/:id/apis 角色分配接口
func (h *RoleController) AssignAPIs(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req RoleAPIAssignReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	if err := h.Svc.AssignRoleAPIs(id, req.APIIDs); err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, nil)
}
