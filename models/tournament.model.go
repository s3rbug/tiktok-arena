package models

import "github.com/google/uuid"

type Tournament struct {
	ID      *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name    string     `gorm:"not null"`
	UserID  *uuid.UUID
	User    User `gorm:"foreignKey:UserID"`
	Tiktoks []TiktokReference
}

type CreateTournament struct {
	Name    string
	UserID  string
	Tiktoks []CreateTiktok
}

func GetAllowedTournamentSize() map[int32]bool {
	return map[int32]bool{
		8:  true,
		16: true,
		32: true,
		64: true,
	}
}

func CheckIfAllowedTournamentSize(sizeToCheck int32) bool {
	return GetAllowedTournamentSize()[sizeToCheck]
}
