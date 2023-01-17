package jwtutil

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"train-tiktok/common/tool"
	"train-tiktok/service/identity/internal/svc"
)

func GenerateJwt(l *svc.ServiceContext, userid int64, username string) (string, error) {

	uInfo := UserInfo{
		UserId:   userid,
		Username: username,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, JwtClaims{
		ID:       generateJti(uInfo),
		ExpireAt: time.Now().Add(24 * 30 * time.Hour).Unix(),
		IssuedAt: time.Now().Unix(),
		Issuer:   "train-tiktok",
		Username: uInfo.Username,
		UserId:   uInfo.UserId,
	}).SignedString(l.JwtSigningKey)
}

// CheckPermission
// @description check user permission
func CheckPermission(l *svc.ServiceContext, _jwt string) (UserInfo, error) {
	JwtClaims, err := JwtDecode(l, _jwt)
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{
		UserId:   JwtClaims.UserId,
		Username: JwtClaims.Username,
	}, nil
}

// JwtDecode decode jwt
func JwtDecode(l *svc.ServiceContext, _jwt string) (JwtClaims, error) {
	token, err := jwt.ParseWithClaims(_jwt, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return l.JwtSigningKey, nil
	})
	if err != nil {
		return JwtClaims{}, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return *claims, nil
	}

	return JwtClaims{}, errors.New("jwt token error")
}

// generateJti 创建jwt唯一标识, 用于后续吊销凭证
func generateJti(user UserInfo) string {
	JtiJson, _ := json.Marshal(map[string]string{
		"username": user.Username,
		"randStr":  tool.RandStr(32),
		"time":     time.Now().Format("150405"),
	})
	_jti := tool.Md5(string(JtiJson))
	return _jti
}

type JwtClaims struct {
	ID       string `json:"jti,omitempty"`
	ExpireAt int64  `json:"exp,omitempty"`
	IssuedAt int64  `json:"iat,omitempty"`
	Issuer   string `json:"iss,omitempty"`
	Username string `json:"username"`
	UserId   int64  `json:"userId"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
}
