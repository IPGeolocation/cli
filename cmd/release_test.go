package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReleaseGuardsKeepExistingArtifacts(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("release assembly runs on Unix")
	}
	source, err := os.ReadFile("../build-release.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ compiler, tag, want string }{
		{"go1.24.5", "", "Use Go"}, {"go1.26.9", "", "Use Go"}, {"go1.27.2", "", "Use Go"},
		{"go1.28rc1", "", "Use Go"}, {"go1.26.10", "v2.1.0", "Release tag must match"},
		{"go1.26.10", "v2.1.1", "checks-started"}, {"go1.27.3", "", "checks-started"}, {"go1.28.0", "", "checks-started"},
	} {
		t.Run(tc.compiler+tc.tag, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.Mkdir(filepath.Join(dir, "dist"), 0700); err != nil {
				t.Fatal(err)
			}
			for name, content := range map[string]string{
				"build-release.sh": string(source), "dist/keep": "existing archive",
				"go":          "#!/bin/sh\ncase \"$1\" in\nenv) echo \"$TEST_COMPILER\";;\nrun) echo 'ipgeolocation version 2.1.1';;\n*) echo checks-started >&2; exit 7;;\nesac\n",
				"govulncheck": "#!/bin/sh\nexit 0\n",
			} {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0700); err != nil {
					t.Fatal(err)
				}
			}
			process := exec.Command("bash", filepath.Join(dir, "build-release.sh"))
			process.Env = append(os.Environ(), "PATH="+dir+string(os.PathListSeparator)+os.Getenv("PATH"), "TEST_COMPILER="+tc.compiler, "GITHUB_REF_TYPE=tag", "GITHUB_REF_NAME="+tc.tag)
			if tc.tag == "" {
				process.Env = append(process.Env, "GITHUB_REF_TYPE=branch")
			}
			output, err := process.CombinedOutput()
			if err == nil || !strings.Contains(string(output), tc.want) {
				t.Fatalf("unexpected guard result: %v\n%s", err, output)
			}
			data, err := os.ReadFile(filepath.Join(dir, "dist/keep"))
			if err != nil || string(data) != "existing archive" {
				t.Fatalf("guard destroyed existing artifacts: %q %v", data, err)
			}
		})
	}
}
