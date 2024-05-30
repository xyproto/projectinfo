package projectinfo

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// StripCodeBlockMarkers removes leading and trailing Markdown code block markers (```)
func StripCodeBlockMarkers(input string) string {
	re := regexp.MustCompile(`(?ms)^\\s*\\x60{3}(?:[a-zA-Z]+)?\\n(.*?)\\n\\x60{3}$`)
	if matches := re.FindStringSubmatch(input); len(matches) > 1 {
		return matches[1]
	}
	return input
}

// URLFromGitConfig reads the .git/config file and extracts the repository URL from the "remote origin" section
func URLFromGitConfig(configFilePath string) (string, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inRemoteSection := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[remote \"origin\"]") {
			inRemoteSection = true
		} else if inRemoteSection && strings.Contains(line, "url =") {
			return strings.TrimSpace(strings.Split(line, "=")[1]), nil
		} else if inRemoteSection && line == "" {
			break // Exit the loop on empty line, which typically ends the section
		}
	}
	return "", fmt.Errorf("no URL found in %s", configFilePath)
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
