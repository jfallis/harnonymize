package main

import (
	"errors"
	"harnonymise/pkg/harnonymize"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path, _ := os.Executable()
	path = filepath.Dir(path)
	entries, _ := os.ReadDir(path)
	keywords := ReadBlockContextKeywords(path)

	anon := harnonymize.New()
	anon.BlockContentKeywords = keywords
	for _, entry := range entries {
		har := harnonymize.NewHAR(path, entry.Name())
		if err := anon.Read(har); err != nil {
			if errors.Is(err, harnonymize.ErrNotHARFile) {
				continue
			}

			panic(err)
		}
		anon.Anonymize(har)

		if wErr := anon.Write(har); wErr != nil {
			panic(wErr)
		}
	}
}

func ReadBlockContextKeywords(path string) []string {
	content, _ := os.ReadFile(filepath.Join(path, "block.txt"))
	return strings.Split(strings.ToLower(string(content)), "\n")
}
