package database

import (
	"testing"
)

func TestNew(t *testing.T) {
	db := New()

	if db.DB == nil {
		t.Error("db load failed")
	}
}
