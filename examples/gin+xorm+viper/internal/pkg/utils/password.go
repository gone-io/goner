package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
	"log"
)

const (
	saltSize   = 16    // 盐的字节长度
	iterations = 10000 // 迭代次数
	keyLength  = 32    // 生成的哈希长度
)

// 生成随机盐
func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// 使用 PBKDF2 生成密码哈希，自动加盐
func hashPassword(password string, salt []byte) string {
	hash := pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)
	return base64.StdEncoding.EncodeToString(append(salt, hash...)) // 将盐和哈希组合并编码为字符串
}

func GenPasswordHash(password string) (string, error) {
	salt, err := generateSalt(saltSize)
	if err != nil {
		return "", err
	}
	return hashPassword(password, salt), nil
}

// VerifyPassword 验证密码是否正确
func VerifyPassword(storedHash string, password string) bool {
	decodedHash, err := base64.StdEncoding.DecodeString(storedHash)
	if err != nil {
		log.Println("Decode Error:", err)
		return false
	}

	salt := decodedHash[:saltSize]
	hash := decodedHash[saltSize:]

	// 使用同样的盐和密码重新生成哈希
	computedHash := pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha256.New)

	// 使用 subtle.ConstantTimeCompare 防止计时攻击
	return subtle.ConstantTimeCompare(hash, computedHash) == 1
}
