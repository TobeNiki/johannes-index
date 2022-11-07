package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuidObj.String(), nil
}

func Password2hash(password string) string {
	salt := "johannes"
	binary := sha256.Sum256([]byte(password + salt))
	hashedPassword := hex.EncodeToString(binary[:])
	return hashedPassword
}

func IncludeStringSlice(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
func GetHitStringFromSlice(slice []string, target string) string {
	//この関数は上記のIncludeStringSliceを使った後に用いること
	for _, str := range slice {
		if str == target {
			return str
		}
	}
	return ""
}
