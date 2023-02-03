package userutil

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"train-tiktok/common/tool"
	"train-tiktok/service/identity/common/errx"
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
	token, err := jwt.ParseWithClaims(_jwt, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return l.JwtSigningKey, nil
	})
	if err != nil {
		return models.User{}, errx.ErrTokenInvalid
	}

	_jwtClaims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return models.User{}, errx.ErrTokenInvalid
	}

	return models.User{
		ID:       _jwtClaims.UserId,
		Username: _jwtClaims.Username,
	}, nil
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
