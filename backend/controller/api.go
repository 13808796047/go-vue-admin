package controller

import (
	"go-vue-admin/repository"
	"go-vue-admin/service"
	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

// APIController 接口权限（api_permissions 表）
type APIController struct {
	Svc *service.Service
}

func NewAPIController(svc *service.Service) *APIController {
	return &APIController{Svc: svc}
}

// List GET /api/apis 接口权限列表（分页）
func (h *APIController) List(c *gin.Context) {
	var q RoleQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.Error(c, service.CodeInvalid, "参数错误: "+err.Error())
		return
	}
	p := &repository.Paginate{Page: q.Page, PageSize: q.PageSize}
	p.Normalize()
	list, total, err := h.Svc.PageAPIs(p)
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, utils.PageResult{Total: total, Page: p.Page, PageSize: p.PageSize, List: list})
}

// All GET /api/apis/all 全量（用于角色分配接口下拉）
func (h *APIController) All(c *gin.Context) {
	list, err := h.Svc.ListAllAPIs()
	if err != nil {
		respond(c, err)
		return
	}
	utils.Success(c, list)
}
