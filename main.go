package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("usage: ./crawler <URL> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}
	if len(args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseURL, maxConcurrencyStr, maxPagesStr := os.Args[1], os.Args[2], os.Args[3]

	maxConcurrency, err := strconv.Atoi(maxConcurrencyStr)
	if err != nil {
		fmt.Printf("Can't convert '%s' to number: %v\n", maxConcurrencyStr, err)
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(maxPagesStr)
	if err != nil {
		fmt.Printf("Can't convert '%s' to number: %v\n", maxPagesStr, err)
		os.Exit(1)
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Can't parse URL %s\n", rawBaseURL)
		os.Exit(1)
	}

	cfg := config{
		maxPages:           maxPages,
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	printReport(cfg.pages, rawBaseURL)
}
