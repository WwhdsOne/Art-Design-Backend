package jwt

import (
	"Art-Design-Backend/pkg/constant"
	"Art-Design-Backend/pkg/redisx"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

// 错误类型定义
var (
	TokenExpired     = errors.New("令牌已过期")
	TokenMalformed   = errors.New("令牌格式错误")
	TokenNotValidYet = errors.New("令牌尚未生效")
	TokenInvalid     = errors.New("令牌无效")
)

// JWT 结构体定义
type JWT struct {
	SigningKey  []byte        // 密钥
	ExpiresTime time.Duration // 过期时间
	Issuer      string        // 签发人
	Audience    string        // 接受者
}

// JWT 实例
var jwtInstance JWT

func NewJWT(initJWT JWT) {
	jwtInstance = initJWT
}

// BaseClaims 基础声明结构体
type BaseClaims struct {
	ID            int64         // 主键 id
	refreshWindow time.Duration // 刷新 token 前多少时间过期
}

// CustomClaims 自定义声明结构体
type CustomClaims struct {
	BaseClaims           // 基础 claims
	jwt.RegisteredClaims // 注册 claims
}

// IsWithinRefreshWindow 判断是否在刷新时间窗口内
func IsWithinRefreshWindow(c *CustomClaims) bool {
	now := time.Now()
	expireTime := c.ExpiresAt.Time
	return now.Add(c.BaseClaims.refreshWindow).After(expireTime)
}

// CreateToken 创建一个 token
func CreateToken(baseClaims BaseClaims) (tokenStr string, err error) {
	// 创建 JWT claims，包含用户 ID
	claim := CreateClaims(baseClaims)
	// 将用户 ID 转换为字符串
	id := strconv.FormatInt(baseClaims.ID, 10)
	// 检查是否存在当前用户的会话
	existToken := redisx.Get(constant.SESSION + id)
	// 如果会话已存在，删除先前对话和 token
	if existToken != "" {
		delKeys := []string{constant.LOGIN + existToken, constant.SESSION + id}
		err = redisx.PipelineDelete(delKeys)
		if err != nil {
			return
		}
	}
	// 创建 JWT 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err = token.SignedString(jwtInstance.SigningKey)
	if err != nil {
		return
	}
	keyVal := [][2]string{
		{constant.LOGIN + tokenStr, id},
		{constant.SESSION + id, tokenStr},
	}
	// 设置会话和 token
	err = redisx.PipelineSet(keyVal, jwtInstance.ExpiresTime)
	return
}

// LogoutToken 注销 token
func LogoutToken(tokenStr string) (err error) {
	claims, err := ParseToken(tokenStr)
	if err != nil {
		return
	}
	// 获取用户 ID
	id := strconv.FormatInt(claims.BaseClaims.ID, 10)
	// 删除 Redis 中的会话信息和 token
	delKeys := []string{constant.SESSION + id, constant.LOGIN + tokenStr}
	err = redisx.PipelineDelete(delKeys)
	if err != nil {
		return
	}
	return err
}

// CreateClaims 创建负载
func CreateClaims(baseClaims BaseClaims) CustomClaims {
	// 设置刷新窗口时间
	baseClaims.refreshWindow = jwtInstance.ExpiresTime / 20
	claims := CustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{jwtInstance.Audience},                      // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1)),                      // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtInstance.ExpiresTime)), // 过期时间 24 小时  配置文件
			Issuer:    jwtInstance.Issuer,                                          // 签名的发行者
		},
	}
	return claims
}

// ParseToken 解析 token
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 使用 jwt.ParseWithClaims 方法解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtInstance.SigningKey, nil // 返回签名密钥
		},
	)

	// 检查解析过程中是否发生错误
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired // 令牌过期
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed // 令牌格式错误
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet // 令牌未生效
		default:
			return nil, TokenInvalid // 其他错误视为无效令牌
		}
	}

	// 确保令牌有效且声明类型正确
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil // 返回自定义声明
	}

	// 如果令牌无效或声明类型不匹配，返回无效令牌错误
	return nil, TokenInvalid
}
