package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID          *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string     `gorm:"not null"`
	Password    string     `gorm:"not null"`
	Tournaments []Tournament
}

type RegisterInput struct {
	Name     string
	Password string
}

type LoginInput struct {
	Name     string
	Password string
}
