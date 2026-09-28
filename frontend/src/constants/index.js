// 后端业务码约定：code === 200 为成功
export const CODE_SUCCESS = 200

// 权限标识（与后端 menus.permKey 对应，前端按钮/菜单据此渲染）
// super_admin 角色持有全部权限，isSuper 时绕过 hasPermission 检查
export const PERMISSION = {
  USER: {
    VIEW: 'sys:user',
    ADD: 'sys:user:add',
    EDIT: 'sys:user:edit',
    DELETE: 'sys:user:delete',
  },
  ROLE: {
    VIEW: 'sys:role',
    ADD: 'sys:role:add',
    DELETE: 'sys:role:delete',
  },
  MENU: {
    VIEW: 'sys:menu',
    ADD: 'sys:menu:add',
    DELETE: 'sys:menu:delete',
  },
}

// 路由页面权限（与后端菜单 type=2 的 permKey 对应）
export const ROUTE_PERM = {
  userManage: 'sys:user',
  roleManage: 'sys:role',
  menuManage: 'sys:menu',
  permission: 'sys:role',
  profile: null,
}

// 通用开关状态
export const STATUS = {
  ENABLED: 1,
  DISABLED: 0,
}

export const STATUS_TEXT = {
  1: '启用',
  0: '停用',
}

export const STATUS_TYPE = {
  1: 'success',
  0: 'info',
}

// 用户状态（1=正常, 0=停用）
export const USER_STATUS = {
  ACTIVE: 1,
  DISABLED: 0,
}

export const USER_STATUS_TEXT = {
  1: '正常',
  0: '停用',
}

export const USER_STATUS_TYPE = {
  1: 'success',
  0: 'info',
}
