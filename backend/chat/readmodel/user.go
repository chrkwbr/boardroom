package readmodel

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Name string
	Icon string
}
