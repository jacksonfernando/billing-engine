package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHexadecimalSHA256Hash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	byteHash := hash.Sum(nil)
	hexaDecimalString := hex.EncodeToString(byteHash)
	return hexaDecimalString
}
