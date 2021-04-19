package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	)

type User struct {
	gorm.Model
	Email string `json:"email" gorm:"not null;unique"`
	Password string `json:"-" gorm:"not null"`
	IsActive bool `json:"isActive" gorm:"default:true"`
	IsLocked bool `json:"isLocked" gorm:"default:false"`
}

type UserPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func PayloadToUser(p UserPayload) User {
	return User{
		Email: p.Email,
		Password: p.Password,
	}
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
/*
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.Password = nil

	return
}
*/
