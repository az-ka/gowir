package util

import "github.com/google/uuid"

func MustNewUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic("critical system error: failed to generate UUID: " + err.Error())
	}
	return id
}
