package main

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func dbCreateUser(user *User) error {
	hashPwd, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashPwd)
	if result := db.Create(user); result.Error != nil {
		return result.Error
	}
	return nil
}
