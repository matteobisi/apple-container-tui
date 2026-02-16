package unit

import (
	"strings"
	"testing"
	"time"

	"container-tui/src/services"
)

func TestImageListParsePerformance100Rows(t *testing.T) {
	var builder strings.Builder
	builder.WriteString("NAME TAG DIGEST\n")
	for i := 0; i < 100; i++ {
		builder.WriteString("repo/image")
		builder.WriteString(string(rune('A' + (i % 26))))
		builder.WriteString(" latest sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n")
	}

	start := time.Now()
	images, err := services.ParseImageList(builder.String())
	duration := time.Since(start)
	if err != nil {
		t.Fatalf("expected no parse error, got %v", err)
	}
	if len(images) != 100 {
		t.Fatalf("expected 100 images, got %d", len(images))
	}
	if duration > 2*time.Second {
		t.Fatalf("expected parsing under 2s, got %s", duration)
	}
}

func TestShellDetectionLatencyBudget(t *testing.T) {
	executor := &shellDetectorExecutor{}
	detector := services.NewShellDetector(executor)

	start := time.Now()
	_, err := detector.DetectShell("latency-test")
	duration := time.Since(start)
	if err != nil {
		t.Fatalf("expected shell detection to succeed, got %v", err)
	}
	if duration > 100*time.Millisecond {
		t.Fatalf("expected shell detection under 100ms in unit context, got %s", duration)
	}
}