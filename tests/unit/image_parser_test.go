package unit

import (
	"testing"

	"container-tui/src/services"
)

func TestParseImageListWellFormed(t *testing.T) {
	output := "NAME TAG DIGEST\nubuntu latest sha256:abc123\n"
	images, err := services.ParseImageList(output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}
	if images[0].Name != "ubuntu" || images[0].Tag != "latest" {
		t.Fatalf("unexpected image: %+v", images[0])
	}
}

func TestParseImageListEmpty(t *testing.T) {
	images, err := services.ParseImageList("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(images) != 0 {
		t.Fatalf("expected no images, got %d", len(images))
	}
}

func TestParseImageListMalformed(t *testing.T) {
	_, err := services.ParseImageList("BROKEN HEADER\nvalue")
	if err == nil {
		t.Fatalf("expected parse error")
	}
}
