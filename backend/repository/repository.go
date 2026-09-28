package repository

import (
	"go-vue-admin/model"

	"gorm.io/gorm"
)

// Base 通用 repository
type Base struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Base {
	return &Base{DB: db}
}

// Paginate 分页参数
type Paginate struct {
	Page     int `form:"page" binding:"required,gt=0"`
	PageSize int `form:"pageSize" binding:"required,gt=0"`
}

func (p *Paginate) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 10
	}
}

// Offset
func (p *Paginate) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// UserRepo
type UserRepo struct{ *Base }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{Base: New(db)} }

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.DB.Preload("Roles").Where("username = ? AND deleted_at IS NULL", username).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepo) GetByID(id uint) (*model.User, error) {
	var u model.User
	if err := r.DB.Preload("Roles").First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Page(p *Paginate, keyword string, roleID uint, status *int8) ([]model.User, int64, error) {
	query := r.DB.Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", like, like, like)
	}
	if roleID > 0 {
		query = query.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Where("user_roles.role_id = ?", roleID)
	}
	if status != nil {
		query = query.Where("users.status = ?", *status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []model.User
	err := query.Preload("Roles").
		Order("users.id ASC").
		Offset(p.Offset()).Limit(p.PageSize).
		Find(&list).Error
	return list, total, err
}

func (r *UserRepo) Create(u *model.User, roleIDs []uint) error {
	u.Roles = make([]model.Role, 0, len(roleIDs))
	for _, id := range roleIDs {
		u.Roles = append(u.Roles, model.Role{ID: id})
	}
	return r.DB.Create(u).Error
}

func (r *UserRepo) Update(u *model.User, roleIDs []uint) error {
	roles := make([]model.Role, 0, len(roleIDs))
	for _, id := range roleIDs {
		roles = append(roles, model.Role{ID: id})
	}
	return r.DB.Model(u).Association("Roles").Replace(&roles)
}

// UpdateProfile 更新资料
func (r *UserRepo) UpdateProfile(u *model.User) error {
	return r.DB.Model(u).Updates(map[string]any{
		"nickname": u.Nickname,
		"email":    u.Email,
		"phone":    u.Phone,
		"avatar":   u.Avatar,
	}).Error
}

func (r *UserRepo) UpdatePassword(id uint, hash string) error {
	return r.DB.Model(&model.User{}).Where("id = ?", id).Update("password", hash).Error
}

func (r *UserRepo) UpdateStatus(id uint, status int8) error {
	return r.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRepo) SoftDelete(id uint) error {
	return r.DB.Delete(&model.User{}, id).Error
}

// RoleRepo
type RoleRepo struct{ *Base }

func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{Base: New(db)} }

func (r *RoleRepo) Page(p *Paginate, keyword string) ([]model.Role, int64, error) {
	query := r.DB.Model(&model.Role{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []model.Role
	err := query.Preload("Menus").Preload("APIs").
		Order("roles.sort ASC, roles.id ASC").
		Offset(p.Offset()).Limit(p.PageSize).
		Find(&list).Error
	return list, total, err
}

func (r *RoleRepo) GetByID(id uint) (*model.Role, error) {
	var ro model.Role
	if err := r.DB.Preload("Menus").Preload("APIs").First(&ro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &ro, nil
}

func (r *RoleRepo) Create(ro *model.Role) error {
	return r.DB.Create(ro).Error
}

func (r *RoleRepo) Update(ro *model.Role) error {
	return r.DB.Model(ro).Updates(map[string]any{
		"name":   ro.Name,
		"sort":   ro.Sort,
		"status": ro.Status,
		"remark": ro.Remark,
	}).Error
}

func (r *RoleRepo) Delete(id uint) error {
	// 解绑
	if err := r.DB.Exec("DELETE FROM user_roles WHERE role_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.DB.Exec("DELETE FROM role_menus WHERE role_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.DB.Exec("DELETE FROM role_apis WHERE role_id = ?", id).Error; err != nil {
		return err
	}
	return r.DB.Delete(&model.Role{}, id).Error
}

func (r *RoleRepo) AssignMenus(roleID uint, menuIDs []uint) error {
	menus := make([]model.Menu, 0, len(menuIDs))
	for _, id := range menuIDs {
		menus = append(menus, model.Menu{ID: id})
	}
	return r.DB.Model(&model.Role{ID: roleID}).
		Association("Menus").Replace(&menus)
}

func (r *RoleRepo) AssignAPIs(roleID uint, apiIDs []uint) error {
	apis := make([]model.API, 0, len(apiIDs))
	for _, id := range apiIDs {
		apis = append(apis, model.API{ID: id})
	}
	return r.DB.Model(&model.Role{ID: roleID}).
		Association("APIs").Replace(&apis)
}

// MenuRepo
type MenuRepo struct{ *Base }

func NewMenuRepo(db *gorm.DB) *MenuRepo { return &MenuRepo{Base: New(db)} }

// Tree 菜单树
func (r *MenuRepo) Tree() ([]model.Menu, error) {
	var all []model.Menu
	if err := r.DB.Order("sort ASC, id ASC").Find(&all).Error; err != nil {
		return nil, err
	}
	children := make(map[uint][]model.Menu)
	var roots []model.Menu
	for _, m := range all {
		if m.ParentID == 0 {
			roots = append(roots, m)
		} else {
			children[m.ParentID] = append(children[m.ParentID], m)
		}
	}
	var build func(id uint) []model.Menu
	build = func(id uint) []model.Menu {
		kids := children[id]
		for i := range kids {
			kids[i].Children = build(kids[i].ID)
		}
		return kids
	}
	for i := range roots {
		roots[i].Children = build(roots[i].ID)
	}
	return roots, nil
}
func (r *MenuRepo) GetByID(id uint) (*model.Menu, error) {
	var mn model.Menu
	if err := r.DB.First(&mn, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mn, nil
}

func (r *MenuRepo) Create(mn *model.Menu) error {
	return r.DB.Create(mn).Error
}

func (r *MenuRepo) Update(mn *model.Menu) error {
	return r.DB.Model(mn).Updates(map[string]any{
		"parent_id": mn.ParentID,
		"name":      mn.Name,
		"perm_key":  mn.PermKey,
		"type":      mn.Type,
		"icon":      mn.Icon,
		"url":       mn.URL,
		"component": mn.Component,
		"sort":      mn.Sort,
		"status":    mn.Status,
	}).Error
}

func (r *MenuRepo) Delete(id uint) error {
	// 有子菜单不允许删除
	var childCount int64
	if err := r.DB.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return gorm.ErrRecordNotFound // 约定：由 service 转为业务错误
	}
	return r.DB.Delete(&model.Menu{}, id).Error
}

func (r *MenuRepo) List() ([]model.Menu, error) {
	var list []model.Menu
	err := r.DB.Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

// APIRepo
type APIRepo struct{ *Base }

func NewAPIRepo(db *gorm.DB) *APIRepo { return &APIRepo{Base: New(db)} }

func (r *APIRepo) Page(p *Paginate) ([]model.API, int64, error) {
	var total int64
	if err := r.DB.Model(&model.API{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.API
	err := r.DB.Order("id ASC").Offset(p.Offset()).Limit(p.PageSize).Find(&list).Error
	return list, total, err
}

func (r *APIRepo) ListAll() ([]model.API, error) {
	var list []model.API
	err := r.DB.Order("id ASC").Find(&list).Error
	return list, err
}
