package entity

import (
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
)

func GenerateBaseId() string {
	return shortuuid.New()
}

func GenerateUuid() string {
	return uuid.NewString()
}
