package models

import "github.com/google/uuid"

type Tiktok struct {
	ID           *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	TournamentID uint
	Tournament   Tournament `gorm:"foreignKey:TournamentID"`
	URL          string     `gorm:"not null"`
	Wins         int
	AvgPoints    float64
}
