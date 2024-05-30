package projectinfo

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// FileInfo represents information about a file in the project, including its content-related attributes
type FileInfo struct {
	Path         string   `json:"path"`
	Language     string   `json:"language"`
	LastModified string   `json:"last_modified,omitempty"`
	Contents     string   `json:"contents,omitempty"`
	LineCount    int      `json:"line_count,omitempty"`
	TokenCount   int      `json:"token_count"`
	Contributors []string `json:"contributors"`
}

// CollectFiles walks through a directory recursively and collects files that have the right extensions
func CollectFiles(dir string, ignores map[string]struct{}, alsoDocOrConf, alsoContributors, verbose bool) ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if verbose {
			fmt.Printf("Visiting: %s\n", path)
		}
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Continue to the next file
		}
		if ShouldSkip(path, ignores) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil // Skip file
		}
		if !d.IsDir() && (RecognizedExtension(path, alsoDocOrConf) || RecognizedFilename(path, alsoDocOrConf)) {
			ext := filepath.Ext(path)
			language := LanguageFromExtension(ext)
			if language != "Unknown" {
				fi, err := os.Stat(path)
				if err != nil {
					log.Printf("Error getting file info for %s: %v\n", path, err)
					return nil // Continue to the next file
				}
				content, err := os.ReadFile(path)
				if err != nil {
					log.Printf("Error reading file %s: %v\n", path, err)
					return nil // Continue to the next file
				}
				utf8Content, err := ConvertToUTF8(content)
				if err != nil {
					log.Printf("Error converting file %s to UTF-8: %v\n", path, err)
					return nil // Continue to the next file
				}
				stringContent := string(utf8Content)
				lineCount, _ := CountLines(stringContent)
				fileInfo := FileInfo{
					Path:         path,
					Language:     language,
					LineCount:    lineCount,
					LastModified: fi.ModTime().Format("2006-01-02 15:04:05"),
					Contents:     stringContent,
				}
				if alsoContributors {
					fileInfo.Contributors = maybeGitContributorsForFile(path)
				}
				files = append(files, fileInfo)
			}
		}
		return nil
	})
	return files, err
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
