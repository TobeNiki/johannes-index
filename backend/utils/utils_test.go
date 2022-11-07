package utils

import (
	"fmt"
	"testing"
)

//動作確認的UnitTest
func TestGenerateUUID(t *testing.T) {
	uuid, _ := GenerateUUID()
	fmt.Println(uuid)
	if uuid == "" {
		t.Error("miss")
	}
}
func TestPassword2Hash(t *testing.T) {
	userInputPass := "test"
	hashedPassword := Password2hash(userInputPass)
	fmt.Println(hashedPassword)
	if hashedPassword == "" {
		t.Error("miss")
	}
}
