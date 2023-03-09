package models

import "github.com/google/uuid"

type FiberMessage struct {
	Message string
}

type UserInfo struct {
	ID   *uuid.UUID
	Name string
}

type AuthDetails struct {
	Username string
	Token    string
}
