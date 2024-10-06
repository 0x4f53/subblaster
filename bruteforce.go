package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/0x4f53/textsubs"
)

func bruteforce() {
	err := os.MkdirAll("output", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		os.Exit(-1)
	}

	files, err := os.ReadDir(batchCache)
	if err != nil {
		fmt.Println("Error reading files:", err)
		os.Exit(-1)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		entries, err := readFile2(batchCache + file.Name())
		if err != nil {
			fmt.Println("Error reading file:", file.Name(), err)
			continue
		}

		var wg sync.WaitGroup
		workerChan := make(chan struct{}, maxWorkers)

		for _, entry := range entries {
			for _, domain := range domains {
				wg.Add(1)
				go func(subdomain, domain string) {
					defer wg.Done()
					workerChan <- struct{}{}
					defer func() { <-workerChan }()

					outputFilePath := filepath.Join("output", domain+".txt")
					if *pairs {
						outputFilePath = filepath.Join("output", domain+".json")
					}

					outputFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE, 0644)
					if err != nil {
						fmt.Println("Error opening output file:", outputFilePath, err)
						return
					}
					defer outputFile.Close()

					checkAndLogSubdomain(subdomain, domain, outputFile)

				}(entry, domain)
			}
		}

		wg.Wait()

		err = os.Remove(batchCache + file.Name())
		if err != nil {
			fmt.Println("Error deleting file:", file.Name(), err)
		} else {
			fmt.Println("[✓] Processed and deleted file ", file.Name())
		}
	}

	fmt.Println("Finished processing all files and subdomains.")
}

func readFile2(filePath string) ([]string, error) {
	var entries []string

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func checkAndLogSubdomain(subdomain, domain string, outputFile *os.File) {
	fullDomain := fmt.Sprintf("%s.%s", subdomain, domain)
	_, err := net.LookupHost(fullDomain)
	if err == nil {
		fmt.Println(" - " + fullDomain)

		// Check if the subdomain already exists in the file
		exists, readErr := stringInFile(fullDomain, outputFile)
		if readErr != nil {
			fmt.Println("Error reading from output file:", readErr)
			return
		}

		if !exists {
			if *pairs {
				pairedSubdomains, _ := textsubs.SubdomainAndDomainPair(fullDomain, false, false)
				for _, item := range pairedSubdomains {
					jsonOutput, _ := json.Marshal(item)
					_, writeErr := outputFile.WriteString(string(jsonOutput) + "\n")
					if writeErr != nil {
						fmt.Println("Error writing to output file:", writeErr)
					}
				}
			} else {
				_, writeErr := outputFile.WriteString(fullDomain + "\n")
				if writeErr != nil {
					fmt.Println("Error writing to output file:", writeErr)
				}
			}
		}
	}
}
