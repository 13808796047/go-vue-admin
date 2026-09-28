package controller

import (
	"go-vue-admin/model"
	"go-vue-admin/service"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

// MenuController
type MenuController struct {
	Svc *service.Service
}

func NewMenuController(svc *service.Service) *MenuController {
	return &MenuController{Svc: svc}
}

// Tree GET /api/menus 菜单树
func (h *MenuController) Tree(c *gin.Context) {
	tree, err := h.Svc.MenuTree()
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, tree)
}

// All GET /api/menus/all 平铺列表（用于权限树）
func (h *MenuController) All(c *gin.Context) {
	list, err := h.Svc.MenuListAll()
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, list)
}

// Get GET /api/menus/:id
func (h *MenuController) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	mn, err := h.Svc.GetMenu(id)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, mn)
}

// Create POST /api/menus
func (h *MenuController) Create(c *gin.Context) {
	var req MenuCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	mn := &model.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		PermKey:   req.PermKey,
		Type:      req.Type,
		Icon:      req.Icon,
		URL:       req.URL,
		Component: req.Component,
		Sort:      req.Sort,
		Status:    req.Status,
	}
	if mn.Type == 0 {
		mn.Type = 2
	}
	if mn.Status == 0 {
		mn.Status = 1
	}
	res, err := h.Svc.CreateMenu(mn)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, res)
}

// Update PUT /api/menus/:id
func (h *MenuController) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req MenuUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	mn := &model.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		PermKey:   req.PermKey,
		Type:      req.Type,
		Icon:      req.Icon,
		URL:       req.URL,
		Component: req.Component,
		Sort:      req.Sort,
		Status:    req.Status,
	}
	res, err := h.Svc.UpdateMenu(id, mn)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, res)
}

// Delete DELETE /api/menus/:id
func (h *MenuController) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.Svc.DeleteMenu(id); err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, nil)
}
