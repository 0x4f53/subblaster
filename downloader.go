package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func save(url string, wg *sync.WaitGroup, p *mpb.Progress, outputDir string) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	fileName := filepath.Join(outputDir, filepath.Base(url))
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file for %s: %v\n", url, err)
		return
	}
	defer file.Close()

	// Progressbar
	bar := p.AddBar(0,
		mpb.BarStyle("╢▌▌░╟"),
		mpb.PrependDecorators(decor.Name(" - "+filepath.Base(url)), decor.CountersKiloByte(" [%.2f / %.2f]")),
		mpb.AppendDecorators(decor.EwmaSpeed(decor.UnitKB, " %.2f", 60)),
	)

	progressReader := bar.ProxyReader(resp.Body)
	defer progressReader.Close()

	_, err = io.Copy(file, progressReader)
	if err != nil {
		fmt.Printf("Failed to write data from %s: %v\n", url, err)
		return
	}
}

func fetch(urls []string) {
	var wg sync.WaitGroup

	p := mpb.New(mpb.WithWaitGroup(&wg))

	for _, url := range urls {
		wg.Add(1)
		go save(url, &wg, p, wordlistCache)
	}

	wg.Wait()
}
