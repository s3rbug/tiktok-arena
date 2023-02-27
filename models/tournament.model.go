package models

import "github.com/google/uuid"

type Tournament struct {
	ID      *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name    string     `gorm:"not null"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	Tiktoks []Tiktok
}
