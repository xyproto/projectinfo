package projectinfo

import (
	"encoding/json"
	"fmt"
	"log"
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

var maxTokensPerChunk = 250000 // This constant defines the maximum token count per chunk of JSON output to avoid excessive data in one batch.

func SetMaxTokensPerChunk(n int) {
	maxTokensPerChunk = n
}

func MaxTokensPerChunk() int {
	return maxTokensPerChunk
}

func New(dir string, outputWarnings bool) (ProjectInfo, error) {
	projectName, err := ReadProjectName(dir)
	if err != nil && outputWarnings {
		projectName = "Untitled"
		log.Printf("could not find project name, using %q: %v\n", projectName, err)
	}

	repoURL, err := URLFromGitConfig(filepath.Join(dir, ".git", "config"))
	if err != nil && outputWarnings {
		log.Printf("could not find git url from git config: %v\n", err)
	}

	ignores, err := LoadIgnorePatterns(dir, ".ignore", ".gitignore")
	if err != nil && outputWarnings {
		log.Printf("could not read .ignore and/or .gitignore: %v\n", err)
	}

	var alsoDocAndConf bool
	const alsoGitContributors = true
	sourceFiles, err := CollectFiles(dir, ignores, alsoDocAndConf, alsoGitContributors)
	if err != nil && outputWarnings {
		log.Printf("could not collect source files: %v\n", err)
	}

	alsoDocAndConf = true
	confAndDocFiles, err := CollectFiles(dir, ignores, alsoDocAndConf, alsoGitContributors)
	if err != nil && outputWarnings {
		log.Printf("could not collect documentation and config files: %v\n", err)
	}

	contributors, err := GitContributors(dir)
	if err != nil && outputWarnings {
		log.Printf("could not collect contributor names from git: %v\n", err)
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

func (project *ProjectInfo) AllFiles() []FileInfo {
	return append(project.SourceFiles, project.ConfAndDocFiles...)
}

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
		if currentTokenCount+file.TokenCount > maxTokensPerChunk {
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
