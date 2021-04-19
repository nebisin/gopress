package models

import (
	"errors"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
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

func (u User) Validate(action string) error {
	validate := validator.New()

	switch strings.ToLower(action) {
	case "register":
		if err := validate.Var(u.Email, "required,email"); err != nil {
			return errors.New("you have to provide a valid email")
		}

		if len(u.Password) < 8 {
			return errors.New("password must be at least 8 characters")
		}
	}

	return nil
}
