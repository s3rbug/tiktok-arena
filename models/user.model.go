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

type UserAuthDetails struct {
	ID       string
	Username string
	Token    string
}
