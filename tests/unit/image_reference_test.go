package unit

import (
	"testing"

	"container-tui/src/models"
)

func TestImageReferenceValidate(t *testing.T) {
	cases := []struct {
		name      string
		reference models.ImageReference
		wantErr   bool
	}{
		{
			name: "valid repository only",
			reference: models.ImageReference{
				Repository: "library/alpine",
			},
			wantErr: false,
		},
		{
			name: "valid with tag",
			reference: models.ImageReference{
				Repository: "nginx",
				Tag:        "latest",
			},
			wantErr: false,
		},
		{
			name:      "missing repository",
			reference: models.ImageReference{},
			wantErr:   true,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.reference.Validate()
			if testCase.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
