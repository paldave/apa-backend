package entity

import "github.com/lithammer/shortuuid/v3"

func GenerateBaseId() string {
	return shortuuid.New()
}
