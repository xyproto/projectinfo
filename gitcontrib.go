package projectinfo

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// TODO: Use a Go module instead of the git command

// maybeGitContriburorsForFile returns a slice of contributors for a given file or directory
func maybeGitContributorsForFile(path string) []string {
	dir := path
	if !isDir(path) {
		dir = filepath.Dir(path)
	}
	cmd := exec.Command("git", "shortlog", "-sn", "--all", "--no-merges", path)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var contributors []string
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if name := ParseContributor(line); name != "" {
			contributors = append(contributors, name)
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}
	}
	return contributors
}

// Contributors uses the Git command line to fetch a list of contributors sorted by the number of commits
func GitContributors(path string) ([]string, error) {
	dir := path
	if !isDir(path) {
		dir = filepath.Dir(path)
	}
	cmd := exec.Command("git", "shortlog", "-sn", "--all", "--no-merges")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return []string{}, fmt.Errorf("failed to execute git command: %v", err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var contributors []string
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if name := ParseContributor(line); name != "" {
			contributors = append(contributors, name)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning git output: %v", err)
	}
	return contributors, nil
}

// parseContributor extracts the name of the contributor from a line of git shortlog output.
func ParseContributor(line string) string {
	parts := strings.SplitN(line, "\t", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
