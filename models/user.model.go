package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string     `gorm:"not null"`
	Password string     `gorm:"not null"`
}

type RegisterInput struct {
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

type LoginInput struct {
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

type UserInfo struct {
	ID   *uuid.UUID
	Name string
}

type UserAuthDetails struct {
	Username string
	Token    string
}
