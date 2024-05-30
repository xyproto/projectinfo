package projectinfo

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadIgnorePatterns loads patterns from specified filenames to ignore during file operations
func LoadIgnorePatterns(filenames ...string) (map[string]struct{}, error) {
	ignores := make(map[string]struct{})
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			continue // skip files that cannot be read
		}
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				ignores[line] = struct{}{}
			}
		}
	}
	// Add common ignores typically found in projects
	commonIgnores := []string{"vendor", "test", "tmp", "backup", "node_modules", "target", ".mvn", ".gradle", ".git"}
	for _, dir := range commonIgnores {
		ignores[dir] = struct{}{}
	}
	return ignores, nil
}

// ShouldSkip determines if a directory should be skipped based on ignore patterns
func ShouldSkip(path string, ignores map[string]struct{}) bool {
	for ignore := range ignores {
		if matched, _ := filepath.Match(ignore, filepath.Base(path)); matched {
			return true
		}
		if strings.HasPrefix(path, ignore+string(filepath.Separator)) {
			return true
		}
	}
	return false
}
