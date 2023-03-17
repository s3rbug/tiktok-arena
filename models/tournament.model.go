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
	Size    int            `validate:"gte=4,lte=64"`
	Tiktoks []CreateTiktok `validate:"required"`
}

type Bracket struct {
	CountMatches int
	Rounds       []Round
}

type Round struct {
	Round   int
	Matches []Match
}

type Option interface {
	isOption() bool
}

type MatchOption struct {
	MatchID string
}

func (m MatchOption) isOption() bool {
	return true
}

type TiktokOption struct {
	TiktokURL string
}

func (m TiktokOption) isOption() bool {
	return true
}

type Match struct {
	MatchID      string
	FirstOption  interface{}
	SecondOption interface{}
}

type ContestPayload struct {
	ContestType string `validate:"required"`
}

const (
	SingleElimination = "single_elimination"
	KingOfTheHill     = "king_of_the_hill"
)

func GetAllowedTournamentType() map[string]bool {
	return map[string]bool{
		SingleElimination: true,
		KingOfTheHill:     true,
	}
}

func CheckIfAllowedTournamentType(tournamentType string) bool {
	return GetAllowedTournamentType()[tournamentType]
}
