package projectinfo

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// FileInfo represents information about a file in the project, including its content-related attributes
type FileInfo struct {
	Path         string `json:"path"`
	Language     string `json:"language"`
	LastModified string `json:"last_modified,omitempty"`
	Contents     string `json:"contents,omitempty"`
	LineCount    int    `json:"line_count,omitempty"`
	TokenCount   int    `json:"token_count"`
}

// CollectFiles walks through a directory recursively and collects files that have the right extensions
func CollectFiles(dir string, ignores map[string]struct{}, catalogFlag bool) ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && ShouldSkip(path, ignores) {
			return fs.SkipDir
		}
		if !d.IsDir() && RecognizedExtension(path, catalogFlag) {
			ext := filepath.Ext(path)
			language := LanguageFromExtension(ext)
			if language != "Unknown" {
				fileInfo, err := os.Stat(path)
				if err != nil {
					return err
				}
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				utf8Content, err := ConvertToUTF8(content)
				if err != nil {
					return err // Handle files that cannot be converted to UTF-8
				}
				stringContent := string(utf8Content)
				lineCount, _ := CountLines(stringContent)

				if relativePath, err := filepath.Rel(dir, path); err == nil { // success
					path = relativePath
				}

				files = append(files, FileInfo{
					Path:         path,
					Language:     language,
					LineCount:    lineCount,
					LastModified: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
					Contents:     stringContent,
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

// ConvertToUTF8 attempts to convert a byte slice to UTF-8 encoding, managing non-UTF8 encoded parts.
func ConvertToUTF8(b []byte) ([]byte, error) {
	if utf8.Valid(b) {
		return b, nil // Return as is if the byte slice is already valid UTF-8
	}
	validUTF8 := bytes.Buffer{}
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			fmt.Println("Encountered invalid UTF-8 sequence, skipping")
			b = b[1:]
			continue
		}
		validUTF8.WriteRune(r)
		b = b[size:]
	}
	return validUTF8.Bytes(), nil
}

// FindFileName returns the first file that has a filename that matches the given string
func FindFileName(files []FileInfo, searchString string) FileInfo {
	for _, fileInfo := range files {
		if strings.Contains(fileInfo.Path, searchString) {
			return fileInfo
		}
	}
	return FileInfo{}
}
