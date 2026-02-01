package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 错误类型定义
var (
	// ErrTokenExpired 令牌已过期
	ErrTokenExpired = errors.New("令牌已过期")
	// ErrTokenMalformed 令牌格式错误
	ErrTokenMalformed = errors.New("令牌格式错误")
	// ErrTokenNotValidYet 令牌尚未生效
	ErrTokenNotValidYet = errors.New("令牌尚未生效")
	// ErrTokenInvalid 令牌无效
	ErrTokenInvalid = errors.New("令牌无效")
)

// JWT 结构体定义
// JWT 提供了JWT令牌的生成和验证功能
type JWT struct {
	SigningKey  []byte        // 密钥
	Issuer      string        // 签发人
	Audience    string        // 接受者
	ExpiresTime time.Duration // 过期时间
}

// BaseClaims 基础声明结构体
// BaseClaims 定义了JWT令牌的基础声明
type BaseClaims struct {
	UserID int64 // 用户主键ID
}

// NewBaseClaims 创建基础声明
// NewBaseClaims 根据用户ID创建BaseClaims实例
func NewBaseClaims(userID int64) BaseClaims {
	return BaseClaims{
		UserID: userID,
	}
}

// CustomClaims 自定义声明结构体
// CustomClaims 定义了JWT令牌的自定义声明
type CustomClaims struct {
	BaseClaims           // 基础 claims
	jwt.RegisteredClaims // 注册 claims
}

// createClaims 创建负载
func (j *JWT) createClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{j.Audience},                      // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1)),            // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresTime)), // 过期时间 24 小时  配置文件
			Issuer:    j.Issuer,                                          // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个 token
func (j *JWT) CreateToken(baseClaims BaseClaims) (tokenStr string, err error) {
	claim := j.createClaims(baseClaims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err = token.SignedString(j.SigningKey)
	return
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotValidYet
		default:
			return nil, ErrTokenInvalid
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}
