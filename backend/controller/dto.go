package controller

// ===== 认证 =====
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ===== 用户 =====
type UserCreateReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

type UserUpdateReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

type UserStatusReq struct {
	Status uint `json:"status"`
}

type UserQuery struct {
	Keyword  string `form:"keyword"`
	RoleID   uint   `form:"roleId"`
	Status   *int8  `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// ===== 角色 =====
type RoleCreateReq struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type RoleUpdateReq struct {
	Name   string `json:"name" binding:"required"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status"`
	Remark string `json:"remark"`
}

type RoleMenuAssignReq struct {
	MenuIDs []uint `json:"menuIds"`
}

type RoleAPIAssignReq struct {
	APIIDs []uint `json:"apiIds"`
}

type RoleQuery struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// ===== 菜单 =====
type MenuCreateReq struct {
	ParentID  uint   `json:"parentId"`
	Name      string `json:"name" binding:"required"`
	PermKey   string `json:"permKey"`
	Type      int8   `json:"type"`
	Icon      string `json:"icon"`
	URL       string `json:"url"`
	Component string `json:"component"`
	Sort      int    `json:"sort"`
	Status    int8   `json:"status"`
}

type MenuUpdateReq struct {
	ParentID  uint   `json:"parentId"`
	Name      string `json:"name" binding:"required"`
	PermKey   string `json:"permKey"`
	Type      int8   `json:"type"`
	Icon      string `json:"icon"`
	URL       string `json:"url"`
	Component string `json:"component"`
	Sort      int    `json:"sort"`
	Status    int8   `json:"status"`
}

// ===== 个人中心 =====
type ProfileUpdateReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

type PasswordChangeReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// ===== 通用 =====
type AssignReq struct {
	IDs []uint `json:"ids"`
}
