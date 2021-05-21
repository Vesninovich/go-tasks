package uuid_test

import (
	"strings"
	"testing"

	"github.com/Vesninovich/go-tasks/book-store/common/uuid"
)

func TestString(t *testing.T) {
	id := uuid.New()
	str := id.String()
	fromStr, err := uuid.FromString(str)
	if err != nil {
		t.Errorf("Failed to create UUID from string created from freshly created UUID: %s", err)
	}
	if id != fromStr {
		t.Errorf("Wrong UUID after stringifying/destringifying:\n\tsource %s\n\tresult %s", id, fromStr)
	}
}

func TestFromInvalidString(t *testing.T) {
	t.Run("Short string", func(t *testing.T) {
		_, err := uuid.FromString(uuid.New().String()[:15])
		if err == nil {
			t.Error("Expected to get error on short string")
		}
	})
	t.Run("Long string", func(t *testing.T) {
		_, err := uuid.FromString(uuid.New().String() + "0")
		if err == nil {
			t.Error("Expected to get error on long string")
		}
	})
	t.Run("Missing -", func(t *testing.T) {
		_, err := uuid.FromString(strings.Replace(uuid.New().String(), "-", "0", 1))
		if err == nil {
			t.Error("Expected to get error on incorrect char where - was expected")
		}
	})
	t.Run("Non-hex character", func(t *testing.T) {
		str := []byte(uuid.New().String())
		str[14] = 'x'
		_, err := uuid.FromString(string(str))
		if err == nil {
			t.Error("Expected to get error on non-hex character")
		}
	})
}
