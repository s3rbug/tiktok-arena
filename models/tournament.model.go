package models

import (
	"github.com/google/uuid"
)

type Tournament struct {
	ID     *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name   string     `gorm:"not null"`
	Size   int        `gorm:"not null"`
	UserID *uuid.UUID
	User   User `gorm:"foreignKey:UserID"`
}

type CreateTournament struct {
	Name    string         `validate:"required"`
	Size    int            `validate:"required"`
	Tiktoks []CreateTiktok `validate:"required"`
}

func GetAllowedTournamentSize() map[int]bool {
	return map[int]bool{
		8:  true,
		16: true,
		32: true,
		64: true,
	}
}

func CheckIfAllowedTournamentSize(sizeToCheck int) bool {
	return GetAllowedTournamentSize()[sizeToCheck]
}
