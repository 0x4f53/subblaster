package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var batchCache = ".batch_cache/"

var batchSize = 500

func BatchFiles(srcDir, destDir string, filenames []string) error {
	for _, filename := range filenames {
		err := processBatchfile(srcDir, destDir, filename)
		if err != nil {
			return fmt.Errorf("error processing file %s: %w", filename, err)
		}
	}
	return nil
}

func processBatchfile(srcDir, destDir, filename string) error {
	srcFilePath := filepath.Join(srcDir, filename)
	file, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", srcFilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0, batchSize)
	batchNum := 1

	for scanner.Scan() {
		lines = append(lines, scanner.Text())

		if len(lines) >= batchSize {
			err := writeBatchToFile(destDir, filename, batchNum, lines)
			if err != nil {
				return err
			}
			batchNum++
			lines = lines[:0] // Reset the slice for the next batch
		}
	}

	if len(lines) > 0 {
		err := writeBatchToFile(destDir, filename, batchNum, lines)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filename, err)
	}

	return nil
}

func writeBatchToFile(destDir, filename string, batchNum int, lines []string) error {
	batchFilename := fmt.Sprintf("%s-%d.txt", filename, batchNum)
	destFilePath := filepath.Join(destDir, batchFilename)

	outFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", destFilePath, err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file %s: %w", destFilePath, err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing file %s: %w", destFilePath, err)
	}

	//fmt.Printf("Wrote batch %d to %s\n", batchNum, destFilePath)
	return nil
}

func batcher() {
	srcDir := wordlistCache
	destDir := batchCache
	filenames, _ := listFiles(wordlistCache)

	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Fatalf("could not create destination directory: %v", err)
	}

	err = BatchFiles(srcDir, destDir, filenames)
	if err != nil {
		log.Fatalf("error batching files: %v", err)
	}
}
