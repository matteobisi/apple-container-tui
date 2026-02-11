package unit

import (
	"testing"

	"container-tui/src/services"
)

func TestTypeToConfirmValidation(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		input    string
		wantOK   bool
	}{
		{"matches", "prod-db", "prod-db", true},
		{"case mismatch", "prod-db", "PROD-DB", false},
		{"empty input", "prod-db", "", false},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			ok := services.IsExactMatch(testCase.expected, testCase.input)
			if ok != testCase.wantOK {
				t.Fatalf("expected %v, got %v", testCase.wantOK, ok)
			}
		})
	}
}
