package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

type CustomClaims struct {
	BaseClaims           // 基础claims
	jwt.RegisteredClaims // 注册claims
}

type BaseClaims struct {
	ID int64 // 主键id
}

type JWT struct {
	SigningKey  []byte        //密钥
	ExpiresTime time.Duration //过期时间
	Issuer      string        //签发人
	Audience    string        //接受者
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateClaims 创建负载
func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{j.Audience},                      // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1)),            // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresTime)), // 过期时间 24小时  配置文件
			Issuer:    j.Issuer,                                          // 签名的发行者
		},
	}
	return claims
}

// ParseToken 解析 JWT 令牌并返回自定义声明（CustomClaims）
// 参数:
//   - tokenString: 要解析的 JWT 令牌字符串
//
// 返回值:
//   - *request.CustomClaims: 解析成功时返回的自定义声明
//   - error: 解析过程中发生的错误，可能的错误包括令牌过期、格式错误、未生效或无效
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 使用 jwt.ParseWithClaims 方法解析令牌
	// 第一个参数是令牌字符串
	// 第二个参数是自定义声明的类型
	// 第三个参数是一个回调函数，用于提供签名密钥
	token, err := jwt.ParseWithClaims(tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 返回签名密钥
			return j.SigningKey, nil
		},
	)

	// 如果解析过程中发生错误，根据错误类型返回相应的错误信息
	if err != nil {
		// 令牌过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		}
		// 令牌格式错误
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, TokenMalformed
		}
		// 令牌未生效
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenNotValidYet
		}
		// 其他错误，视为无效令牌
		return nil, TokenInvalid
	}

	// 如果令牌解析成功
	if token != nil {
		// 尝试将令牌的声明部分转换为自定义声明类型
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			// 如果转换成功且令牌有效，返回自定义声明
			return claims, nil
		}
		// 如果转换失败或令牌无效，返回无效令牌错误
		return nil, TokenInvalid
	} else {
		// 如果令牌解析失败，返回无效令牌错误
		return nil, TokenInvalid
	}
}
