package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户（软删除）
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password    string         `gorm:"size:128;not null" json:"-"`
	Nickname    string         `gorm:"size:64" json:"nickname"`
	Email       string         `gorm:"size:128" json:"email"`
	Phone       string         `gorm:"size:32" json:"phone"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Status      int8           `gorm:"default:1" json:"status"` // 1启用 0禁用
	LastLoginAt *time.Time     `json:"lastLoginAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Roles       []Role         `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles"`
}

func (User) TableName() string { return "users" }

// Role 角色
type Role struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code      string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int8           `gorm:"default:1" json:"status"`
	Remark    string         `gorm:"size:255" json:"remark"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Menus     []Menu         `gorm:"many2many:role_menus;joinForeignKey:RoleID;joinReferences:MenuID" json:"menus,omitempty"`
	APIs      []API          `gorm:"many2many:role_apis;joinForeignKey:RoleID;joinReferences:ApiID" json:"apis,omitempty"`
}

func (Role) TableName() string { return "roles" }

// Menu 菜单/权限
type Menu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ParentID  uint           `gorm:"default:0;index" json:"parentId"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	PermKey   string         `gorm:"size:128;uniqueIndex" json:"permKey"` // 权限标识,如 sys:user:add
	Type      int8           `gorm:"default:1" json:"type"`               // 1目录 2菜单 3按钮/权限
	Icon      string         `gorm:"size:64" json:"icon"`
	URL       string         `gorm:"size:255" json:"url"`
	Component string         `gorm:"size:255" json:"component"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int8           `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Children  []Menu         `gorm:"-" json:"children"`
}

func (Menu) TableName() string { return "menus" }

// API 接口权限
type API struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Path      string         `gorm:"size:255;not null" json:"path"`
	Method    string         `gorm:"size:16;not null" json:"method"`
	Name      string         `gorm:"size:128" json:"name"`
	Desc      string         `gorm:"size:255" json:"desc"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (API) TableName() string { return "api_permissions" }

// 中间表（many2many）
type UserRoles struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (UserRoles) TableName() string { return "user_roles" }

type RoleMenus struct {
	RoleID uint `gorm:"primaryKey"`
	MenuID uint `gorm:"primaryKey"`
}

func (RoleMenus) TableName() string { return "role_menus" }

type RoleApis struct {
	RoleID uint `gorm:"primaryKey"`
	ApiID  uint `gorm:"primaryKey"`
}

func (RoleApis) TableName() string { return "role_apis" }
