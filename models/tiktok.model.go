package models

import "github.com/google/uuid"

type Tiktok struct {
	TournamentID *uuid.UUID  `gorm:"not null;primaryKey"`
	Tournament   *Tournament `gorm:"foreignKey:TournamentID"`
	URL          string      `gorm:"not null;primaryKey"`
	Wins         int
	AvgPoints    float64
	TimesPlayed  int
}

type CreateTiktok struct {
	URL string `validate:"required"`
}
