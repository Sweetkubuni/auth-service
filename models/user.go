package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

const (
	scryptN      = 16384
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 64
)

func (user *User) HashPassword(password string) error {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	dk, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return err
	}

	user.Password = base64.StdEncoding.EncodeToString(dk)
	user.Salt = base64.StdEncoding.EncodeToString(salt)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return err
	}

	dk, err := scrypt.Key([]byte(providedPassword), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return err
	}

	encodedHash := base64.StdEncoding.EncodeToString(dk)
	if user.Password != encodedHash {
		return fmt.Errorf("invalid password")
	}
	return nil
}
