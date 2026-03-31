package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ImageReference represents a container image reference.
type ImageReference struct {
	Registry   string
	Repository string
	Tag        string
	Digest     string
}

// Validate ensures required fields are present.
func (r ImageReference) Validate() error {
	if strings.TrimSpace(r.Repository) == "" {
		return errors.New("repository is required")
	}
	return nil
}

// String formats the reference for CLI usage.
func (r ImageReference) String() string {
	reference := strings.TrimSpace(r.Repository)
	if r.Registry != "" {
		reference = strings.TrimSpace(r.Registry) + "/" + reference
	}
	if r.Digest != "" {
		return reference + "@" + strings.TrimSpace(r.Digest)
	}
	if r.Tag != "" {
		return reference + ":" + strings.TrimSpace(r.Tag)
	}
	return reference
}

// BuildExportImageReference creates a deterministic temporary image reference for export workflows.
func BuildExportImageReference(containerName, containerID string, now time.Time) string {
	slug := ExportNameSlug(containerName, containerID)
	return fmt.Sprintf("actui-export/%s:%s", slug, now.UTC().Format("20060102-150405"))
}

// BuildExportArchiveName creates the generated OCI archive filename for an export workflow.
func BuildExportArchiveName(containerName, containerID string, now time.Time) string {
	slug := ExportNameSlug(containerName, containerID)
	return fmt.Sprintf("%s-%s.oci.tar", slug, now.UTC().Format("20060102-150405"))
}

// ExportNameSlug normalizes a container name or id for generated references and file names.
func ExportNameSlug(containerName, containerID string) string {
	base := strings.TrimSpace(containerName)
	if base == "" {
		base = strings.TrimSpace(containerID)
	}
	if base == "" {
		base = "container"
	}
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		" ", "-",
		"\t", "-",
		"\n", "-",
		"_", "-",
	)
	slug := strings.ToLower(replacer.Replace(base))
	parts := strings.FieldsFunc(slug, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-')
	})
	if len(parts) == 0 {
		return "container"
	}
	joined := strings.Join(parts, "-")
	joined = strings.Trim(joined, "-")
	if joined == "" {
		return "container"
	}
	return joined
}
