package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKubernetesAndVersionDocumentation(t *testing.T) {
	readme := readDocumentation(t, "README.md")
	menuMap := readDocumentation(t, "docs", "ai-menu-map.md")
	userGuide := readDocumentation(t, "docs", "user-guide.md")
	binaryGuide := readDocumentation(t, "docs", "binary-build-automation.md")

	for _, required := range []string{"`k`", "Kubernetes", "actui version dev", "release-version artifact"} {
		if !strings.Contains(readme, required) {
			t.Errorf("README missing %q", required)
		}
	}
	for _, required := range []string{"KubernetesClusterList", "KubernetesClusterSubmenu", "kubernetes_builders.go", "type-to-confirm"} {
		if !strings.Contains(menuMap, required) {
			t.Errorf("AI menu map missing %q", required)
		}
	}
	for _, required := range []string{"container k8s list", "container k8s create", "container k8s start", "container k8s delete", "container k8s load-image", "container k8s write-config", "--dry-run", "exact cluster name"} {
		if !strings.Contains(userGuide, required) {
			t.Errorf("user guide missing %q", required)
		}
	}
	for _, required := range []string{"release-version", "before compilation", "main.version=<tag without v>", "without recomputing a tag"} {
		if !strings.Contains(binaryGuide, required) {
			t.Errorf("binary automation guide missing %q", required)
		}
	}
}

func readDocumentation(t *testing.T, pathElements ...string) string {
	t.Helper()
	path := filepath.Join(append([]string{"..", ".."}, pathElements...)...)
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(contents)
}
