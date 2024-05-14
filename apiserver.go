package projectinfo

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// apiIndicators defines keywords or phrases that suggest a file is part of an API server.
var apiIndicators = map[string][]string{
	".go":   {"ListenAndServe", "HandleFunc", "echo.New()", "gin.Default()", "mux.NewRouter()"},
	".c":    {"#include <mongoose.h>"},
	".cpp":  {"#include <cpprest/http_listener.h>"},
	".hs":   {"Scotty", "Spock"},
	".js":   {"express()", "require('koa')", "fastify()"},
	".ts":   {"@nestjs", "express()", "import 'koa'"},
	".kt":   {"ktor.application"},
	".java": {"@RestController", "@GetMapping", "@PostMapping"},
	// Add other languages and frameworks as needed
}

// PossiblyAPIServer tries to detect if the source code in the given directory is likely to be an API server.
func PossiblyAPIServer(dir string) bool {
	hasAPIIndicators := false
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if indicators, ok := apiIndicators[ext]; ok {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				for _, indicator := range indicators {
					if strings.Contains(line, indicator) {
						hasAPIIndicators = true
						break // Break early since we found an indicator
					}
				}
				if hasAPIIndicators {
					break // Break early from file scanning as we have confirmed it's likely an API file
				}
			}
		}
		return nil
	})
	if err != nil {
		return true // Assume true if there are errors to err on the side of caution
	}
	return hasAPIIndicators
}
