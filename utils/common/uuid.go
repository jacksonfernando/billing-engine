package common

import (
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

func UUIDTOBytes(input uuid.UUID) ([]byte, error) {
	hexStrUUID := strings.ReplaceAll(input.String(), "-", "")
	byteUUID, err := hex.DecodeString(hexStrUUID)
	if err != nil {
		return nil, err
	}
	return byteUUID, nil
}

func GenerateNewUUIDV7() ([]byte, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	res, err := UUIDTOBytes(uuid)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func StringToUUID(input string) (*uuid.UUID, error) {
	id, err := uuid.Parse(input)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
