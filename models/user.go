package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"unique;not null" json:"name"`
	Email string `gorm:"unique;not null" json:"email"`
}
