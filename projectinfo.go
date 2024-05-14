package projectinfo

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// ProjectInfo holds information about the entire project, useful for generating documentation or other reports.
type ProjectInfo struct {
	Name            string     `json:"name"`
	RepoURL         string     `json:"repositoryURL"`
	SourceFiles     []FileInfo `json:"sourceFiles"`
	ConfAndDocFiles []FileInfo `json:"confAndDocFiles"`
	Type            string     `json:"type"`
	Contributors    string     `json:"contributors"`
	APIServer       bool       `json:"apiServer"`
}

func New(dir string) (ProjectInfo, error) {
	projectName, err := ReadProjectName(dir)
	if err != nil {
		// TODO: warn if verbose
		projectName = "Untitled"
	}

	repoURL, err := URLFromGitConfig(filepath.Join(dir, ".git", "config"))
	if err != nil {
		// TODO: warn if verbose
	}

	ignores, err := LoadIgnorePatterns(dir, ".ignore", ".gitignore")
	if err != nil {
		// TODO: warn if verbose
	}

	sourceFiles, err := CollectFiles(dir, ignores, false)
	if err != nil {
		// TODO: warn if verbose
	}

	confAndDocFiles, err := CollectFiles(dir, ignores, true)
	if err != nil {
		// TODO: warn if verbose
	}

	contributors, err := GitContributors(dir)
	if err != nil {
		// TODO: warn if verbose
	}

	apiServer := PossiblyAPIServer(dir)

	return ProjectInfo{
		Name:            projectName,
		RepoURL:         repoURL,
		SourceFiles:     sourceFiles,
		ConfAndDocFiles: confAndDocFiles,
		Type:            DetectProjectType(sourceFiles),
		Contributors:    strings.Join(contributors, ", "),
		APIServer:       apiServer,
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
