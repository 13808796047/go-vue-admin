package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword bcrypt 加密
func HashPassword(plain string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), 10)
	return string(h), err
}

// CheckPassword 校验密码
func CheckPassword(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
