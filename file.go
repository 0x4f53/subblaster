package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var wordlists = "lists.yaml"

var wordlistCache = ".cached_wordlists/"

var combinedFilename = "combined.txt"

func processFile(filePath string, linesChan chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(line) && !strings.Contains(line, " ") {
			line = strings.ToLower(line)
			if line != "" {
				linesChan <- line
			}
		}
	}

	return scanner.Err()
}

func writeToFile(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
	}
	return writer.Flush()
}

func readFile(filePath string, lineChannel chan<- string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineChannel <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func downloadSecLists() {
	if _, err := os.Stat(wordlistCache); os.IsNotExist(err) {
		err := os.Mkdir(wordlistCache, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	secListSources := readListsFile(wordlists)

	fetch(secListSources)

}

func readListsFile(filename string) []string {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var list []string
	err = decoder.Decode(&list)
	if err != nil {
		log.Fatalf("Error decoding YAML: %v", err)
	}

	return list

}

func cacheExists() bool {
	_, err := os.Stat(wordlistCache)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func countLinesInDirectory(directory string) (int, error) {
	totalLines := 0

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			lines, err := countLinesInFile(path)
			if err != nil {
				return err
			}
			totalLines += lines
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to count lines: %w", err)
	}

	return totalLines, nil
}

func countLinesInFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return lineCount, nil
}

func listFiles(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(directory, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, err
}

type FileWriter struct {
	filename string
	mu       sync.Mutex
}

func NewFileWriter(filename string) *FileWriter {
	return &FileWriter{filename: filename}
}

func (fw *FileWriter) WriteToFile(data string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	file, err := os.OpenFile(fw.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	return err
}

func stringInFile(subdomain string, outputFile *os.File) (bool, error) {
	scanner := bufio.NewScanner(outputFile)

	// Reset file pointer to the beginning
	if _, err := outputFile.Seek(0, 0); err != nil {
		return false, err
	}

	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == subdomain {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func downloadAndValidateWordlists() {
	fmt.Println("[↓] Downloading wordlists mentioned in " + wordlists)
	downloadSecLists()
	fmt.Println("\n[⟳] Validating downloaded wordlists...")
	removeDuplicatesFromSecLists()
	fmt.Println("\n[✓] Processing complete!")
	lines, _ := countLinesInDirectory(wordlistCache)
	fmt.Println("\n[+] Wordlist cache count: " + strconv.Itoa(lines) + " items")
}
