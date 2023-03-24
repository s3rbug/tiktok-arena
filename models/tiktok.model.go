package models

import "github.com/google/uuid"

type Tiktok struct {
	TournamentID *uuid.UUID  `gorm:"not null;primaryKey;default:null"`
	Tournament   *Tournament `gorm:"foreignKey:TournamentID"`
	Name         string      `gorm:"not null;default:null"`
	URL          string      `gorm:"not null;primaryKey;default:null"`
	Wins         int
}

type CreateTiktok struct {
	Name string `validate:"required"`
	URL  string `validate:"required"`
}
