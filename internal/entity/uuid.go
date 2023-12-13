package entity

import "github.com/google/uuid"

type UUIDGenerator interface {
	Generate() string
}

type UUIDGeneratorImpl struct{}

func (g UUIDGeneratorImpl) Generate() uuid.UUID {
	return uuid.New()
}
