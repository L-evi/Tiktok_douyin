package userutil

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"train-tiktok/common/tool"
	"train-tiktok/service/identity/internal/svc"
	"train-tiktok/service/identity/models"
)

func GenerateJwt(l *svc.ServiceContext, userid int64, username string) (string, error) {

	uInfo := models.User{
		ID:       userid,
		Username: username,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, JwtClaims{
		ID:       generateJti(uInfo),
		ExpireAt: time.Now().Add(24 * 30 * time.Hour).Unix(),
		IssuedAt: time.Now().Unix(),
		Issuer:   "train-tiktok",
		Username: uInfo.Username,
		UserId:   uInfo.ID,
	}).SignedString(l.JwtSigningKey)
}

// CheckPermission
// @description check user permission
func CheckPermission(l *svc.ServiceContext, _jwt string) (models.User, error) {
	JwtClaims, err := JwtDecode(l, _jwt)
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:       JwtClaims.UserId,
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
func generateJti(user models.User) string {
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
