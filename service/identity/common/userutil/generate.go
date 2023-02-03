package userutil

import (
	"strconv"
	"strings"
	"train-tiktok/common/tool"
)

// GenerateNickname generates a random nickname for new User
func GenerateNickname() string {
	var nickname strings.Builder

	prefix := "抖声用户"
	nickname.WriteString(prefix)
	nickname.WriteString(strconv.Itoa(tool.RandInt(6)))

	return nickname.String()
}
