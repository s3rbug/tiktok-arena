package models

import (
	"github.com/google/uuid"
)

type Tournament struct {
	ID     *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name   string     `gorm:"not null"`
	Size   int        `gorm:"not null"`
	UserID *uuid.UUID `gorm:"not null"`
	User   *User      `gorm:"foreignKey:UserID"`
}

type CreateTournament struct {
	Name    string         `validate:"required"`
	Size    int            `validate:"required"`
	Tiktoks []CreateTiktok `validate:"required"`
}

type SingleEliminationBracket struct {
	CountMatches int
	Rounds       *[]Round
}

type Round struct {
	Round   int
	Matches []Match
}

type MatchOption struct {
	MatchID string
}

type TiktokOption struct {
	TiktokURL string
}

type Match struct {
	MatchID      string
	FirstOption  interface{}
	SecondOption interface{}
}

type ContestItem struct {
	ID           string
	FirstOption  ContestOption
	SecondOption ContestOption
}

type ContestOption struct {
	OptionID string
	Url      string
}

type ContestPayload struct {
	ContestType string `validate:"required"`
}

func GetAllowedTournamentType() map[string]bool {
	return map[string]bool{
		"single elimination": true,
	}
}

func CheckIfAllowedTournamentType(tournamentType string) bool {
	return GetAllowedTournamentType()[tournamentType]
}

func GetAllowedTournamentSize() map[int]bool {
	return map[int]bool{
		8:  true,
		16: true,
		32: true,
	}
}

func CheckIfAllowedTournamentSize(sizeToCheck int) bool {
	return GetAllowedTournamentSize()[sizeToCheck]
}
