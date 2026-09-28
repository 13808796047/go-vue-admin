package service

import (
	"errors"
	"fmt"
	"time"

	"go-vue-admin/config"
	"go-vue-admin/model"
	"go-vue-admin/repository"
	"go-vue-admin/utils"

	"gorm.io/gorm"
)

// 业务错误
var (
	ErrNotFound      = errors.New("记录不存在")
	ErrAlreadyExists = errors.New("记录已存在")
	ErrNoPermission  = errors.New("无访问权限")
	ErrInvalidCred   = errors.New("用户名或密码错误")
	ErrAccountLocked = errors.New("账号已禁用")
	ErrHasChildren   = errors.New("存在子节点，无法删除")
	ErrOldPassword   = errors.New("原密码不正确")
	ErrProtected     = errors.New("内置数据不可修改或删除")
)

// 错误码常量（响应 code）
const (
	CodeOK       = 200
	CodeInvalid  = 400
	CodeUnauth   = 401
	CodeNoPerm   = 403
	CodeNotFound = 404
	CodeServer   = 500
)

// UserPermissions 用户权限结果
type UserPermissions struct {
	Menus     []model.Menu `json:"menus"`
	APIs      []model.API  `json:"apis"`
	PermKeys  []string     `json:"permKeys"`
	RoleCodes []string     `json:"roleCodes"`
}

// LoginResult 登录返回
type LoginResult struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresIn    int         `json:"expiresIn"`
	User         *model.User `json:"user"`
}

// Service 聚合业务逻辑
type Service struct {
	Users *repository.UserRepo
	Roles *repository.RoleRepo
	Menus *repository.MenuRepo
	APIs  *repository.APIRepo
	Cfg   *config.Config
}

func New(cfg *config.Config, db *gorm.DB) *Service {
	return &Service{
		Users: repository.NewUserRepo(db),
		Roles: repository.NewRoleRepo(db),
		Menus: repository.NewMenuRepo(db),
		APIs:  repository.NewAPIRepo(db),
		Cfg:   cfg,
	}
}

// ===== 认证 =====

// Login 用户登录
func (s *Service) Login(username, password string) (*LoginResult, error) {
	u, err := s.Users.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCred
	}
	if u.Status != 1 {
		return nil, ErrAccountLocked
	}
	if !utils.CheckPassword(password, u.Password) {
		return nil, ErrInvalidCred
	}
	access, err := utils.SignToken(s.Cfg.JWTSecret, &utils.Claims{UserID: u.ID, Username: u.Username, Type: "access"}, s.Cfg.JWTTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := utils.SignToken(s.Cfg.JWTSecret, &utils.Claims{UserID: u.ID, Username: u.Username, Type: "refresh"}, s.Cfg.RefreshTTL)
	if err != nil {
		return nil, err
	}
	if err := s.updateLastLogin(u.ID); err != nil {
		return nil, err
	}
	return &LoginResult{AccessToken: access, RefreshToken: refresh, ExpiresIn: s.Cfg.JWTTTL, User: u}, nil
}

func (s *Service) updateLastLogin(uid uint) error {
	return s.Users.DB.Exec("UPDATE users SET last_login_at = ? WHERE id = ?", time.Now(), uid).Error
}

// Refresh 用 refresh token 换取新的 access token
func (s *Service) Refresh(refreshToken string) (string, error) {
	claims, err := utils.ParseToken(s.Cfg.JWTSecret, refreshToken)
	if err != nil {
		return "", errors.New("refresh token 无效或已过期")
	}
	if claims.Type != "refresh" {
		return "", errors.New("token 类型错误")
	}
	return utils.SignToken(s.Cfg.JWTSecret, &utils.Claims{UserID: claims.UserID, Username: claims.Username, Type: "access"}, s.Cfg.JWTTTL)
}

// ===== 用户 =====

// PageUsers 分页查询用户
func (s *Service) PageUsers(p *repository.Paginate, keyword string, roleID uint, status *int8) ([]model.User, int64, error) {
	list, total, err := s.Users.Page(p, keyword, roleID, status)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		list[i].Password = ""
	}
	return list, total, nil
}

// GetUser
func (s *Service) GetUser(id uint) (*model.User, error) {
	u, err := s.Users.GetByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	return u, nil
}

// CreateUser
func (s *Service) CreateUser(username, nickname, password, email, phone, avatar string, status int8, roleIDs []uint) (*model.User, error) {
	existing, err := s.Users.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadyExists
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Username: username,
		Nickname: nickname,
		Password: hash,
		Email:    email,
		Phone:    phone,
		Avatar:   avatar,
		Status:   status,
	}
	if err := s.Users.Create(u, roleIDs); err != nil {
		return nil, err
	}
	return s.GetUser(u.ID)
}

// UpdateUser
func (s *Service) UpdateUser(id uint, nickname, email, phone, avatar string, status int8, roleIDs []uint) (*model.User, error) {
	u, err := s.Users.GetByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	u.Nickname = nickname
	u.Email = email
	u.Phone = phone
	u.Avatar = avatar
	u.Status = status
	if err := s.Users.Update(u, roleIDs); err != nil {
		return nil, err
	}
	return s.GetUser(id)
}

// SetUserStatus
func (s *Service) SetUserStatus(id, status uint) error {
	if id == s.superUserID() {
		return fmt.Errorf("%w: 不能修改内置超级管理员状态", ErrProtected)
	}
	if err := s.Users.UpdateStatus(id, int8(status)); err != nil {
		return err
	}
	return nil
}

// DeleteUser 软删除（内置超管不可删除）
func (s *Service) DeleteUser(id uint) error {
	if u, _ := s.Users.GetByID(id); u != nil && u.Username == "admin" {
		return fmt.Errorf("%w: 不能删除内置超级管理员", ErrProtected)
	}
	if err := s.Users.SoftDelete(id); err != nil {
		return err
	}
	return nil
}

// ===== 角色 =====

func (s *Service) PageRoles(p *repository.Paginate, keyword string) ([]model.Role, int64, error) {
	return s.Roles.Page(p, keyword)
}

func (s *Service) GetRole(id uint) (*model.Role, error) {
	ro, err := s.Roles.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ro == nil {
		return nil, ErrNotFound
	}
	return ro, nil
}

func (s *Service) CreateRole(name, code string, sort int, status int8, remark string) (*model.Role, error) {
	ro := &model.Role{Name: name, Code: code, Sort: sort, Status: status, Remark: remark}
	if err := s.Roles.Create(ro); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || isUniqueViolation(err) {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}
	return s.GetRole(ro.ID)
}

func (s *Service) UpdateRole(id uint, name string, sort int, status int8, remark string) (*model.Role, error) {
	ro, err := s.Roles.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ro == nil {
		return nil, ErrNotFound
	}
	if ro.Code == "super_admin" {
		return nil, fmt.Errorf("%w: 不能修改内置超级管理员角色", ErrProtected)
	}
	ro.Name = name
	ro.Sort = sort
	ro.Status = status
	ro.Remark = remark
	if err := s.Roles.Update(ro); err != nil {
		return nil, err
	}
	return s.GetRole(id)
}

func (s *Service) DeleteRole(id uint) error {
	ro, err := s.Roles.GetByID(id)
	if err != nil {
		return err
	}
	if ro == nil {
		return ErrNotFound
	}
	if ro.Code == "super_admin" {
		return fmt.Errorf("%w: 不能删除内置超级管理员角色", ErrProtected)
	}
	return s.Roles.Delete(id)
}

// AssignRoleMenus 角色分配菜单
func (s *Service) AssignRoleMenus(roleID uint, menuIDs []uint) error {
	ro, err := s.Roles.GetByID(roleID)
	if err != nil {
		return err
	}
	if ro == nil {
		return ErrNotFound
	}
	return s.Roles.AssignMenus(roleID, menuIDs)
}

// AssignRoleAPIs 角色分配接口
func (s *Service) AssignRoleAPIs(roleID uint, apiIDs []uint) error {
	ro, err := s.Roles.GetByID(roleID)
	if err != nil {
		return err
	}
	if ro == nil {
		return ErrNotFound
	}
	return s.Roles.AssignAPIs(roleID, apiIDs)
}

// ===== 菜单 =====

func (s *Service) MenuListAll() ([]model.Menu, error) {
	return s.Menus.List()
}

func (s *Service) MenuTree() ([]model.Menu, error) {
	return s.Menus.Tree()
}

// ===== 接口权限 =====

func (s *Service) PageAPIs(p *repository.Paginate) ([]model.API, int64, error) {
	return s.APIs.Page(p)
}

func (s *Service) ListAllAPIs() ([]model.API, error) {
	return s.APIs.ListAll()
}

func (s *Service) GetMenu(id uint) (*model.Menu, error) {
	mn, err := s.Menus.GetByID(id)
	if err != nil {
		return nil, err
	}
	if mn == nil {
		return nil, ErrNotFound
	}
	return mn, nil
}

func (s *Service) CreateMenu(m *model.Menu) (*model.Menu, error) {
	if err := s.Menus.Create(m); err != nil {
		return nil, err
	}
	return s.GetMenu(m.ID)
}

func (s *Service) UpdateMenu(id uint, m *model.Menu) (*model.Menu, error) {
	mn, err := s.Menus.GetByID(id)
	if err != nil {
		return nil, err
	}
	if mn == nil {
		return nil, ErrNotFound
	}
	m.ID = id
	if err := s.Menus.Update(m); err != nil {
		return nil, err
	}
	return s.GetMenu(id)
}

func (s *Service) DeleteMenu(id uint) error {
	mn, err := s.Menus.GetByID(id)
	if err != nil {
		return err
	}
	if mn == nil {
		return ErrNotFound
	}
	err = s.Menus.Delete(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrHasChildren
	}
	return err
}

// ===== 权限 =====

// UserPermissions 当前用户菜单树 + 接口权限 + 权限标识 + 角色编码
func (s *Service) UserPermissions(userID uint) (*UserPermissions, error) {
	u, err := s.Users.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}

	// 收集角色
	roleIDs := make([]uint, 0, len(u.Roles))
	roleCodes := make([]string, 0, len(u.Roles))
	isSuper := false
	for _, ro := range u.Roles {
		roleIDs = append(roleIDs, ro.ID)
		roleCodes = append(roleCodes, ro.Code)
		if ro.Code == "super_admin" {
			isSuper = true
		}
	}

	res := &UserPermissions{
		PermKeys:  []string{},
		RoleCodes: roleCodes,
	}

	if isSuper {
		// 超管：全部菜单 + 全部接口
		roots, err := s.Menus.Tree()
		if err != nil {
			return nil, err
		}
		res.Menus = roots
		apis, err := s.APIs.ListAll()
		if err != nil {
			return nil, err
		}
		res.APIs = apis
		for _, a := range apis {
			res.PermKeys = append(res.PermKeys, a.Path+"|"+a.Method)
		}
		return res, nil
	}

	// 普通用户：按角色查菜单与接口
	menuIDs, apiIDs := s.rolePermIDs(roleIDs)
	if len(menuIDs) == 0 && len(apiIDs) == 0 {
		res.Menus = []model.Menu{}
		res.APIs = []model.API{}
		return res, nil
	}

	allMenus, err := s.Menus.List()
	if err != nil {
		return nil, err
	}
	mSet := make(map[uint]bool, len(menuIDs))
	for _, id := range menuIDs {
		mSet[id] = true
	}
	res.Menus = filterTree(allMenus, mSet)

	var apis []model.API
	if len(apiIDs) > 0 {
		apis, err = s.filterAPIsByIDs(apiIDs)
		if err != nil {
			return nil, err
		}
	}
	res.APIs = apis

	// 菜单 permKey + 接口 标识
	permSet := make(map[string]bool)
	var collect func(ms []model.Menu)
	collect = func(ms []model.Menu) {
		for _, m := range ms {
			if m.PermKey != "" {
				permSet[m.PermKey] = true
			}
			collect(m.Children)
		}
	}
	collect(res.Menus)
	for _, a := range apis {
		permSet[a.Path+"|"+a.Method] = true
	}
	for k := range permSet {
		res.PermKeys = append(res.PermKeys, k)
	}
	return res, nil
}

func (s *Service) rolePermIDs(roleIDs []uint) (menuIDs, apiIDs []uint) {
	if len(roleIDs) == 0 {
		return
	}
	roles := make(map[uint]*model.Role, len(roleIDs))
	for _, id := range roleIDs {
		ro, _ := s.Roles.GetByID(id)
		if ro != nil {
			roles[id] = ro
		}
	}
	mSet := make(map[uint]bool)
	aSet := make(map[uint]bool)
	for _, ro := range roles {
		for _, m := range ro.Menus {
			mSet[m.ID] = true
		}
		for _, a := range ro.APIs {
			aSet[a.ID] = true
		}
	}
	for id := range mSet {
		menuIDs = append(menuIDs, id)
	}
	for id := range aSet {
		apiIDs = append(apiIDs, id)
	}
	return
}

func (s *Service) filterAPIsByIDs(ids []uint) ([]model.API, error) {
	if len(ids) == 0 {
		return []model.API{}, nil
	}
	var list []model.API
	err := s.APIs.DB.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func filterTree(all []model.Menu, mSet map[uint]bool) []model.Menu {
	children := make(map[uint][]model.Menu)
	var roots []model.Menu
	for _, m := range all {
		if !mSet[m.ID] {
			continue
		}
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
	return roots
}

// ===== 个人中心 =====

// GetProfile 当前用户详情
func (s *Service) GetProfile(userID uint) (*model.User, error) {
	u, err := s.Users.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	u.Password = ""
	return u, nil
}

// UpdateProfile 更新资料
func (s *Service) UpdateProfile(userID uint, nickname, email, phone, avatar string) (*model.User, error) {
	u, err := s.Users.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	u.Nickname = nickname
	u.Email = email
	u.Phone = phone
	u.Avatar = avatar
	if err := s.Users.UpdateProfile(u); err != nil {
		return nil, err
	}
	return s.GetProfile(userID)
}

// ChangePassword 修改密码
func (s *Service) ChangePassword(userID uint, oldPwd, newPwd string) error {
	u, err := s.Users.GetByID(userID)
	if err != nil {
		return err
	}
	if u == nil {
		return ErrNotFound
	}
	if !utils.CheckPassword(oldPwd, u.Password) {
		return ErrOldPassword
	}
	hash, err := utils.HashPassword(newPwd)
	if err != nil {
		return err
	}
	return s.Users.UpdatePassword(userID, hash)
}

// superUserID 内置 admin 的 ID
func (s *Service) superUserID() uint {
	u, _ := s.Users.GetByUsername("admin")
	if u == nil {
		return 0
	}
	return u.ID
}

// isUniqueViolation SQLite 唯一约束
func isUniqueViolation(err error) bool {
	return err != nil && (containsStr(err.Error(), "UNIQUE constraint failed"))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
