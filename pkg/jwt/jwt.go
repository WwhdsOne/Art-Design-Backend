package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
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
	SigningKey        []byte        // 密钥
	Issuer            string        // 签发人
	Audience          string        // 接受者
	ExpiresTime       time.Duration // 过期时间
	RefreshWindowTime time.Duration // 刷新窗口时间
}

// BaseClaims 基础声明结构体
type BaseClaims struct {
	UserID            int64         // 主键 id
	RefreshWindowTime time.Duration // 刷新窗口时间
}

func NewBaseClaims(userId int64) BaseClaims {
	return BaseClaims{
		UserID: userId,
	}
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
	return now.Add(c.RefreshWindowTime).After(expireTime)
}

// createClaims 创建负载
func (j *JWT) createClaims(baseClaims BaseClaims) CustomClaims {
	// 设置刷新窗口时间
	baseClaims.RefreshWindowTime = j.RefreshWindowTime
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
		func(token *jwt.Token) (interface{}, error) {
			return j.SigningKey, nil
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
