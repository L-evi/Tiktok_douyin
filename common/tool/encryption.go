package tool

import "golang.org/x/crypto/bcrypt"

func EncipherPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func VerifyPassword(password, ciphertext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(ciphertext), []byte(password))
	// 出现错误的时候就是对应不上
	if err != nil {
		return false, err
	}
	return true, err
}
