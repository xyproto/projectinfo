package projectname

import (
	"os"
	"path/filepath"
	"testing"
)

func setupMockFile(dir, filename, content string) error {
	tmpDir := filepath.Join(dir, filename)
	return os.WriteFile(tmpDir, []byte(content), 0644)
}

func TestReadProjectName(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
		want     string
	}{
		{
			name:     "package.json valid",
			filename: "package.json",
			content:  `{"name": "testproject"}`,
			want:     "testproject",
		},
		{
			name:     "pom.xml valid",
			filename: "pom.xml",
			content:  `<project><name>JavaProject</name></project>`,
			want:     "JavaProject",
		},
		{
			name:     "csproj valid",
			filename: "project.csproj",
			content:  `<Project><PropertyGroup><AssemblyName>CSharpProject</AssemblyName></PropertyGroup></Project>`,
			want:     "CSharpProject",
		},
		{
			name:     "build.gradle valid",
			filename: "build.gradle",
			content:  `rootProject.name='GradleProject'`,
			want:     "GradleProject",
		},
		{
			name:     "go.mod valid",
			filename: "go.mod",
			content:  `module github.com/example/goProject`,
			want:     "github.com/example/goProject",
		},
		{
			name:     "Cargo.toml valid",
			filename: "Cargo.toml",
			content:  `name = "RustProject"`,
			want:     "RustProject",
		},
		{
			name:     "setup.py valid",
			filename: "setup.py",
			content:  `name='PythonProject'`,
			want:     "PythonProject",
		},
		{
			name:     "cabal file valid",
			filename: "project.cabal",
			content:  `name: HaskellProject`,
			want:     "HaskellProject",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "testdir")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir) // Clean up after the test.

			if err := setupMockFile(tempDir, tc.filename, tc.content); err != nil {
				t.Fatalf("Failed to setup mock file: %v", err)
			}

			got, err := ReadProjectName(tempDir)
			if err != nil {
				t.Errorf("ReadProjectName() error = %v, want no error", err)
			}
			if got != tc.want {
				t.Errorf("ReadProjectName() got = %v, want %v", got, tc.want)
			}
		})
	}
}
