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
	MenuListRole         = "MENU:LIST:ROLE:"         // 角色的菜单信息
	MenuListRoleTTL      = 1 * time.Hour
)

// 用户权限相关
const (
	UserRoleList    = "USER:ROLE:LIST:" // 用户的角色信息
	UserRoleListTTL = 30 * time.Minute  // 缓存有效期
	RoleInfo        = "ROLE:INFO:"
	RoleInfoTTL     = 1 * time.Hour
	// RoleUserDependencies 角色用户依赖关系，不过期
	RoleUserDependencies = "ROLE:USER:DEPENDENCIES:"
)

// AI模型相关
const (
	AIModelSimpleList    = "AIMODEL:SIMPLE:LIST"
	AIModelSimpleListTTL = 86400 * time.Second
	AIModelInfo          = "AIMODEL:INFO:"
	AIModelInfoTTL       = 86400 * time.Second
)

// RateLimiter 访问频率限制
const (
	RateLimiter = "RATE:LIMITER:"
)

// KeyStats 键前缀统计数据
const (
	KeyStats = "KEY_STATS"
)
