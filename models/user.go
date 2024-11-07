package models

import (
	"fmt"
	"sync/atomic"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `gorm:"unique;not null" json:"name"`
	Email string `gorm:"unique;not null" json:"email"`
}

type FakeUserGenerator struct {
	userCounter atomic.Uint64
}

func (generator *FakeUserGenerator) Generate() *User {
	id := generator.userCounter.Add(1)

	return &User{
		Model: gorm.Model{ID: uint(id)},
		Name:  fmt.Sprintf("testuser_%d", id),
		Email: fmt.Sprintf("testuser_%d@example.com", id),
	}
}
