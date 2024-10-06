package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/0x4f53/textsubs"
)

var maxWorkers = 1000000

func logo() {
	fmt.Println(`
		S U B B L A S T E R  
	(https://github.com/0x4f53/subblaster)
	A fast subdomain bruteforcer in Golang.
	`)
}

var domains []string
var pairs *bool

func setFlags() {

	refresh := flag.Bool("r", false, "Re-download all wordlists mentioned in "+wordlists)
	workers := flag.Int("w", maxWorkers, "Set a custom number of workers to use while bruteforcing (default:"+strconv.Itoa(maxWorkers)+")")
	//terraform := flag.Bool("t", false, "Enable Terraform mode (short)")
	pairs = flag.Bool("p", false, `Write paired outputs, e.g.: {"subdomain": "www.example.com", "domain":"example.com"}`)

	flag.Parse()

	domains = flag.Args()

	if *refresh {
		os.RemoveAll(wordlistCache)
		os.RemoveAll(batchCache)
		fmt.Println("[✓] Deleted all cache data")
		downloadAndValidateWordlists()
		fmt.Println("[✓] Refreshed worldists! Re-run this program without the refresh (-r) flag to continue")
		os.Exit(0)
	}

	if len(domains) == 0 && !*refresh {
		fmt.Println(`Please provide at least one domain
E.g.: subblaster example.com
Type "subblaster -h" for more details`)
		os.Exit(-1)
	}

	// Set custom worker amount if mentioned
	maxWorkers = *workers

}

func validateDomains(domains []string) []string {

	var validDomains []string
	for _, domain := range domains {
		isValid, _ := textsubs.DomainsOnly(domain, false)
		if len(isValid) > 0 {
			validDomains = append(validDomains, isValid[0])
		}
	}
	return validDomains
}

func main() {

	logo()

	setFlags()

	domains = validateDomains(domains)

	if !cacheExists() {
		downloadAndValidateWordlists()
	}

	batchCount, _ := listFiles(batchCache)

	if len(batchCount) <= 1 {
		fmt.Println("\n[⟳] Generating batches for bruteforcing...")
		batcher()
		batchCount, _ = listFiles(batchCache)
		fmt.Println("\n[✓] Batching complete! Generated " + strconv.Itoa(len(batchCount)) + " batches")
	} else {
		fmt.Println("\n[+] Found a previous stopping point. Continuing...")
	}

	fmt.Println("\n[+] Bruteforcing...")
	fmt.Println("")

	bruteforce()

}
