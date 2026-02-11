package unit

import (
	"testing"

	"container-tui/src/models"
)

func TestContainerValidate(t *testing.T) {
	cases := []struct {
		name      string
		container models.Container
		wantErr   bool
	}{
		{
			name: "valid running",
			container: models.Container{
				ID:     "abc123",
				Name:   "web",
				Image:  "nginx:latest",
				Status: models.ContainerStatusRunning,
			},
			wantErr: false,
		},
		{
			name: "missing id",
			container: models.Container{
				Name:   "web",
				Image:  "nginx:latest",
				Status: models.ContainerStatusRunning,
			},
			wantErr: true,
		},
		{
			name: "missing name",
			container: models.Container{
				ID:     "abc123",
				Image:  "nginx:latest",
				Status: models.ContainerStatusRunning,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			container: models.Container{
				ID:     "abc123",
				Name:   "web",
				Image:  "nginx:latest",
				Status: "broken",
			},
			wantErr: true,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.container.Validate()
			if testCase.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
