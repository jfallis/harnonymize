package main

import (
	"errors"
	"fmt"
	"harnonymise/pkg/harnonymize"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}

	path = filepath.ToSlash(path)
	path = path[:strings.LastIndex(path, "/")]
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	keywords := ReadBlockContextKeywords(path)

	anon := harnonymize.New()
	anon.BlockContentKeywords = keywords
	for _, entry := range entries {
		har := harnonymize.NewHAR(path, entry.Name())
		if readErr := anon.Read(har); readErr != nil {
			if errors.Is(readErr, harnonymize.ErrNotHARFile) {
				continue
			}
			panic(readErr)
		}

		anon.Anonymize(har)
		if wErr := anon.Write(har); wErr != nil {
			panic(wErr)
		}
	}
}

func ReadBlockContextKeywords(path string) []string {
	name := fmt.Sprintf("%s/%s", path, "block.txt")
	context, readErr := os.ReadFile(name)
	if readErr != nil {
		return nil
	}

	keywords := strings.Split(strings.ToLower(string(context)), "\n")
	return keywords
}
