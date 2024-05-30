package projectinfo

import (
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

func New(dir string, verbose bool) (ProjectInfo, error) {
	projectName, err := ReadProjectName(dir)
	if err != nil && verbose {
		projectName = "Untitled"
		log.Printf("could not find project name, using %q: %v\n", projectName, err)
	}

	repoURL, err := URLFromGitConfig(filepath.Join(dir, ".git", "config"))
	if err != nil && verbose {
		log.Printf("could not find git url from git config: %v\n", err)
	}

	ignores, err := LoadIgnorePatterns(dir, ".ignore", ".gitignore")
	if err != nil && verbose {
		log.Printf("could not read .ignore and/or .gitignore: %v\n", err)
	}

	var alsoDocAndConf bool
	const alsoGitContributors = true
	sourceFiles, err := CollectFiles(dir, ignores, alsoDocAndConf, alsoGitContributors, verbose)
	if err != nil && verbose {
		log.Printf("could not collect source files: %v\n", err)
	}

	alsoDocAndConf = true
	confAndDocFiles, err := CollectFiles(dir, ignores, alsoDocAndConf, alsoGitContributors, verbose)
	if err != nil && verbose {
		log.Printf("could not collect documentation and config files: %v\n", err)
	}

	contributors, err := GitContributors(dir)
	if err != nil && verbose {
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
	var files []FileInfo
	files = append(files, project.SourceFiles...)
	return append(files, project.ConfAndDocFiles...)
}
