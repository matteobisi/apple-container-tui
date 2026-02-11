package unit

import (
	"os"
	"path/filepath"
	"testing"

	"container-tui/src/models"
)

func TestBuildSourceValidate(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "Containerfile")
	if err := os.WriteFile(filePath, []byte("FROM alpine"), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	cases := []struct {
		name    string
		source  models.BuildSource
		wantErr bool
	}{
		{
			name: "valid source",
			source: models.BuildSource{
				FilePath:         filePath,
				FileType:         models.BuildFileTypeContainerfile,
				WorkingDirectory: tempDir,
			},
			wantErr: false,
		},
		{
			name: "missing file path",
			source: models.BuildSource{
				WorkingDirectory: tempDir,
				FileType:         models.BuildFileTypeContainerfile,
			},
			wantErr: true,
		},
		{
			name: "missing working directory",
			source: models.BuildSource{
				FilePath: filePath,
				FileType: models.BuildFileTypeContainerfile,
			},
			wantErr: true,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.source.Validate()
			if testCase.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
