package pkg

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt Password encryption and decryption
type Encrypt struct{}

// NewEncrypt Instantiating Encrypt
func NewEncrypt() *Encrypt {
	return &Encrypt{}
}

// EncryptPWD Encrypted Password
func (*Encrypt) EncryptPWD(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidatePWD Verify Password
func (*Encrypt) ValidatePWD(hash string, b []byte) bool {
	byteHash := []byte(hash)
	if err := bcrypt.CompareHashAndPassword(byteHash, b); err != nil {
		return false
	}
	return true
}

// RandomStr Generate a random string
func (*Encrypt) RandomStr(n int) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)

}

// UniqueStr Generate a random unique string ID
func (*Encrypt) UniqueStr() (string, error) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil

}
