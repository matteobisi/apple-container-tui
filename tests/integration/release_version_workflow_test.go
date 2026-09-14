package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseVersionWorkflow(t *testing.T) {
	buildWorkflow := readWorkflow(t, "build-binary.yml")
	publishWorkflow := readWorkflow(t, "publish-release.yml")

	computedAt := strings.Index(buildWorkflow, "- name: Compute next version tag")
	builtAt := strings.Index(buildWorkflow, "- name: Build actui")
	if computedAt < 0 || builtAt < 0 || computedAt > builtAt {
		t.Fatal("build workflow must compute the release tag before building actui")
	}
	for _, required := range []string{
		"-ldflags \"-X main.version=${RELEASE_TAG#v}\"",
		"- name: Upload release version artifact",
		"name: release-version",
		"path: release-version.txt",
	} {
		if !strings.Contains(buildWorkflow, required) {
			t.Errorf("build workflow missing %q", required)
		}
	}
	for _, required := range []string{
		"- name: Download release version artifact",
		"name: release-version",
		"- name: Read release version",
		"steps.release-version.outputs.tag_name",
	} {
		if !strings.Contains(publishWorkflow, required) {
			t.Errorf("publish workflow missing %q", required)
		}
	}
	if strings.Contains(publishWorkflow, "- name: Compute next version tag") {
		t.Error("publish workflow must use the build version artifact instead of recomputing a tag")
	}
}

func readWorkflow(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join("..", "..", ".github", "workflows", name)
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(contents)
}
