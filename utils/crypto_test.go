package utils

import (
	"regexp"
	"testing"
)

func TestGenerateUUID_Format(t *testing.T) {
	// Arrange
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	// Act
	uuid := GenerateUUID()

	// Assert
	if len(uuid) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid))
	}
	if !uuidPattern.MatchString(uuid) {
		t.Errorf("UUID format invalid: %s", uuid)
	}
}

func TestGenerateUUID_Uniqueness(t *testing.T) {
	// Act
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()
	uuid3 := GenerateUUID()

	// Assert
	if uuid1 == uuid2 || uuid1 == uuid3 || uuid2 == uuid3 {
		t.Errorf("UUIDs not unique: %s, %s, %s", uuid1, uuid2, uuid3)
	}
}
