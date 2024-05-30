package projectinfo

import (
	"path/filepath"
	"strings"
	"unicode"
)

// RecognizedExtension checks if the file extension is recognized and should be included based on the docAndConf flag
func RecognizedExtension(path string, docAndConf bool) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".rst", ".txt", ".adoc", ".md", ".yml", ".yaml", ".properties": // documentation and configuration
		return docAndConf // return true only if catalogFlag is set to true
	case ".c", ".cc", ".cpp", ".cs", ".go", ".h", ".hpp", ".hs", ".java", ".js", ".jsx", ".kt", ".py", ".rs", ".ts", ".tsx", ".sql": // source code only
		return !docAndConf // return true only if catalogFlag is set to false
	}
	return false
}

// RecognizedFilename checks if the given path is a recgonized filename, like "LICENSE"
func RecognizedFilename(path string, docAndConf bool) bool {
	if !docAndConf {
		return false
	}
	switch strings.ToLower(filepath.Base(path)) {
	case "copying", "license", "notice", "makefile":
		return true
	}
	return false
}

// LanguageFromExtension determines the programming language from the file extension
func LanguageFromExtension(ext string) string {
	switch ext {
	case ".adoc":
		return "ASCIIDoc"
	case ".c":
		return "C"
	case ".cpp", ".cc":
		return "C++"
	case ".cs":
		return "C#"
	case ".go":
		return "Go"
	case ".hpp", ".h":
		return "C/C++ Header"
	case ".hs":
		return "Haskell"
	case ".java":
		return "Java"
	case ".js", ".jsx":
		return "JavaScript"
	case ".kt":
		return "Kotlin"
	case ".md":
		return "Markdown"
	case ".properties":
		return "Properties"
	case ".py":
		return "Python"
	case ".rs":
		return "Rust"
	case ".rst":
		return "reStructuredText"
	case ".sql":
		return "SQL"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".txt":
		return "Plain text"
	case ".xml":
		return "XML"
	case ".yml", ".yaml":
		return "YAML"
	default:
		return "Unknown"
	}
}

// DetectProjectType determines the most common programming language used in the project files to suggest the project's type
func DetectProjectType(files []FileInfo) string {
	languageCount := make(map[string]int)
	for _, file := range files {
		languageCount[file.Language]++
	}
	maxCount := 0
	projectType := "Unrecognized"
	for lang, count := range languageCount {
		if count > maxCount {
			maxCount = count
			projectType = lang
		}
	}
	return projectType
}

// OptimizeCode optimizes the source code by normalizing line breaks, trimming unnecessary whitespace, and reducing blank lines.
func OptimizeCode(source string, ext string) string {
	var (
		normalized       = strings.ReplaceAll(source, "\r\n", "\n") // Normalize line breaks
		lines            = strings.Split(normalized, "\n")
		optimizedLines   []string
		lastLineWasBlank bool
		trimLeft         = true // default to trimming whitespace from both ends
	)

	switch ext {
	case ".py", ".hs": // Python, Haskell might need to preserve significant leading spaces
		trimLeft = false
	}

	for _, line := range lines {
		trimmedLine := line
		if trimLeft {
			trimmedLine = strings.TrimSpace(trimmedLine)
		} else {
			trimmedLine = strings.TrimRightFunc(line, unicode.IsSpace)
		}

		if trimmedLine == "" {
			if !lastLineWasBlank {
				optimizedLines = append(optimizedLines, trimmedLine)
				lastLineWasBlank = true
			}
		} else {
			optimizedLines = append(optimizedLines, trimmedLine)
			lastLineWasBlank = false
		}
	}
	return strings.Join(optimizedLines, "\n")
}
