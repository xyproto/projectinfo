package projectinfo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ProjectInfo holds information about the entire project, useful for generating documentation or other reports.
type ProjectInfo struct {
	Name            string     `json:"name"`
	Repository      string     `json:"repository"`
	SourceFiles     []FileInfo `json:"sourceFiles"`
	ConfAndDocFiles []FileInfo `json:"confAndDocFiles"`
	Type            string     `json:"type"`
	Contributors    string     `json:"contributors"`
	CurrentReadMe   FileInfo   `json:"currentReadMe"`
}

func New(dir string) (ProjectInfo, error) {
	ignores, err := LoadIgnorePatterns(dir, ".ignore", ".gitignore")
	if err != nil {
		// skip
		//return ProjectInfo{}, fmt.Errorf("could not load ignore patterns: %v", err)
	}
	sourceFiles, err := CollectFiles(dir, ignores, false)
	if err != nil {
		return ProjectInfo{}, fmt.Errorf("error walking directory and collecting files: %v", err)
	}
	confAndDocFiles, err := CollectFiles(dir, ignores, true)
	if err != nil {
		// skip
		//return ProjectInfo{}, fmt.Errorf("error walking directory and collecting files: %v", err)
	}
	projectName, err := ReadProjectName(dir)
	if err != nil {
		projectName = "Untitled"
		//fmt.Fprintf(os.Stderr, "Warning: could not determine the project name: %v\n", err)
	}
	contributors, err := GitContributors(dir)
	if err != nil {
		// skip
		//return ProjectInfo{}, fmt.Errorf("error collecting git contributors: %v", err)
	}

	return ProjectInfo{
		Name:            projectName,
		SourceFiles:     sourceFiles,
		ConfAndDocFiles: confAndDocFiles,
		Type:            DetectProjectType(sourceFiles),
		Contributors:    strings.Join(contributors, ", "),
	}, nil
}

var MaxTokensPerChunk = 250000 // This constant defines the maximum token count per chunk of JSON output to avoid excessive data in one batch.

// Chunk breaks down project information into manageable JSON chunks to adhere to token limitations
func (project *ProjectInfo) Chunk(includeSourceFiles, includeConfAndDocFiles bool) ([]string, error) {
	var (
		chunks            []string
		currentChunk      []FileInfo
		currentTokenCount = 0
		files             = []FileInfo{}
	)
	if includeSourceFiles {
		files = append(files, project.SourceFiles...)
	}
	if includeConfAndDocFiles {
		files = append(files, project.ConfAndDocFiles...)
	}
	for _, file := range files {
		file.TokenCount = CountTokens(file.Contents) // Compute token count for each file, assuming this function is defined in utils.go.
		if currentTokenCount+file.TokenCount > MaxTokensPerChunk {
			// Finalize the current chunk and reset counters if the maximum token count is exceeded.
			chunkData, err := json.Marshal(currentChunk)
			if err != nil {
				return nil, fmt.Errorf("error marshaling chunk: %v", err)
			}
			chunks = append(chunks, string(chunkData))
			currentChunk = []FileInfo{} // Reset the current chunk
			currentTokenCount = 0
		}
		// Add the file to the current chunk
		currentChunk = append(currentChunk, file)
		currentTokenCount += file.TokenCount
	}
	// Add the last chunk if it contains any files
	if len(currentChunk) > 0 {
		chunkData, err := json.Marshal(currentChunk)
		if err != nil {
			return nil, fmt.Errorf("error marshaling final chunk: %v", err)
		}
		chunks = append(chunks, string(chunkData))
	}
	return chunks, nil
}
