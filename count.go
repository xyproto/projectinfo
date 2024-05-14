package projectinfo

import (
	"bufio"
	"strings"
	"unicode/utf8"
)

// CountLines counts the number of lines in a string, often used to determine the size of file contents
func CountLines(content string) (int, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

// CountTokens estimates the number of tokens in a string based on UTF-8 rune count
func CountTokens(input string) int {
	runeCount := utf8.RuneCountInString(input)
	return (runeCount + 3) / 4 // A rough estimate of token count, suitable for simple use cases
}
