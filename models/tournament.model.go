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

type Bracket struct {
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

func SingleElimination() string {
	return "single_elimination"
}

func KingOfTheHill() string {
	return "king_of_the_hill"
}

func SwissSystem() string {
	return "swiss_system"
}

func DoubleElimination() string {
	return "double_elimination"
}

func GetAllowedTournamentType() map[string]bool {
	return map[string]bool{
		SingleElimination(): true,
		KingOfTheHill():     true,
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
