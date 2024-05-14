package projectinfo

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type PomProject struct {
	XMLName xml.Name `xml:"project"`
	Name    string   `xml:"name"`
}

type CsProj struct {
	XMLName        xml.Name        `xml:"Project"`
	PropertyGroups []PropertyGroup `xml:"PropertyGroup"`
}

type PropertyGroup struct {
	ProjectName string `xml:"AssemblyName"`
}

// ReadProjectName tries to deduce the project name from common configuration files.
func ReadProjectName(dir string) (string, error) {
	checkFunctions := []func(string) (string, error){
		readFromPackageJSON,
		readFromPomXML,
		readFromGradle,
		readFromGoMod,
		readFromCargoToml,
		readFromSetupPy,
		readFromCabal,
		readFromCsProj,
	}

	for _, fn := range checkFunctions {
		if name, err := fn(dir); err == nil {
			return name, nil
		}
	}

	return "", os.ErrNotExist
}

// Define functions to read from each type of configuration file.
func readFromPackageJSON(dir string) (string, error) {
	path := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}
	if name, ok := config["name"].(string); ok {
		return name, nil
	}
	return "", os.ErrNotExist
}

func readFromPomXML(dir string) (string, error) {
	path := filepath.Join(dir, "pom.xml")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var pom PomProject
	if err := xml.Unmarshal(data, &pom); err != nil {
		return "", err
	}
	return pom.Name, nil
}

func readFromGradle(dir string) (string, error) {
	path := filepath.Join(dir, "build.gradle")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "rootProject.name") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				return strings.Trim(strings.TrimSpace(parts[1]), "'"), nil
			}
		}
	}
	return "", os.ErrNotExist
}

func readFromGoMod(dir string) (string, error) {
	path := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[1], nil
			}
		}
	}
	return "", os.ErrNotExist
}

func readFromCargoToml(dir string) (string, error) {
	path := filepath.Join(dir, "Cargo.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "name =") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				return strings.TrimSpace(strings.Trim(parts[1], " \"")), nil
			}
		}
	}
	return "", os.ErrNotExist

}

func readFromSetupPy(dir string) (string, error) {
	path := filepath.Join(dir, "setup.py")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "name=") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				return strings.TrimSpace(strings.Trim(parts[1], ", '\"")), nil
			}
		}
	}
	return "", os.ErrNotExist
}

func readFromCabal(dir string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.cabal"))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", errors.New("could not find *.cabal")
	}
	path := matches[0]
	if len(path) == 0 {
		return "", os.ErrNotExist
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "name:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", os.ErrNotExist
}

func readFromCsProj(dir string) (string, error) {
	path := filepath.Join(dir, "*.csproj")
	matches, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", errors.New("could not find *.csproj")
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		return "", err
	}
	var proj CsProj
	if err := xml.Unmarshal(data, &proj); err != nil {
		return "", err
	}
	if proj.PropertyGroups != nil && len(proj.PropertyGroups) > 0 {
		return proj.PropertyGroups[0].ProjectName, nil
	}
	return "", os.ErrNotExist
}
