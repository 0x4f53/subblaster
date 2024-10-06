package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var globalStringSet = make(map[string]bool)
var mu sync.Mutex

var validSubdomainsRegex *regexp.Regexp

func deduplicate(filename string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	var uniqueStrings []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.ToLower(line)
		if line != "" && validSubdomainsRegex.MatchString(line) {
			mu.Lock()
			if !globalStringSet[line] {
				globalStringSet[line] = true
				uniqueStrings = append(uniqueStrings, line)
			}
			mu.Unlock()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		return
	}

	err = os.WriteFile(filename, []byte(strings.Join(uniqueStrings, "\n")), 0644)
	if err != nil {
		fmt.Printf("Failed to write back to file %s: %v\n", filename, err)
	}
}

func processDirectory(dir string) {
	var wg sync.WaitGroup

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			wg.Add(1)
			go deduplicate(path, &wg)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory %s: %v\n", dir, err)
	}

	wg.Wait()
}

func removeDuplicatesFromSecLists() {
	validSubdomainsRegex, _ = regexp.Compile(`^[\w0-9._-]+$`)
	processDirectory(wordlistCache)
}
