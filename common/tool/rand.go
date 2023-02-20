package tool

import (
	"math/rand"
	"strings"
	"time"
)

const strList string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStr generates a random string with the specified length.
func RandStr(length int) string {
	var result strings.Builder

	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result.WriteByte(strList[rand.Intn(len(strList))])
	}
	return result.String()
}

// RandInt generates a random int with the specified range.
func RandInt(length int) int {
	var code int
	for i := 0; i < length; i++ {
		code += rand.Intn(10)
	}
	return code
}
