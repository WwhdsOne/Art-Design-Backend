package rediskey

import "time"

// 登录相关
// 时限相关内容写在jwt中
// 此处不重复声明
const (
	LOGIN   = "AUTH:LOGIN:"   // token -> userID
	SESSION = "AUTH:SESSION:" // userID -> token
)

// 菜单缓存相关
const (
	// MenuRoleDependencies 菜单角色依赖关系，不过期
	MenuRoleDependencies = "MENU:ROLE:DEPENDENCIES:" // roleID -> dependency tree
	MenuListRole         = "MENU:LIST:ROLE:"         // roleID -> menu list
	MenuListRoleTTL      = 1 * time.Hour
)

// 用户权限相关
const (
	UserRoleList         = "USER:ROLE:LIST:" // userID -> roleID list
	UserRoleListTTL      = 30 * time.Minute  // 缓存有效期
	RoleInfo             = "ROLE:INFO:"
	RoleInfoTTL          = 1 * time.Hour
	RoleUserDependencies = "ROLE:USER:DEPENDENCIES:"
)
