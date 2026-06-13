package domain

import "github.com/google/uuid"

type Room struct {
	ID   uuid.UUID
	Name string
}
