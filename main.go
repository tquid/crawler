package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

const concurrencyLimit = 5

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseURL := os.Args[1]
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Can't parse URL %s\n", rawBaseURL)
		os.Exit(1)
	}

	cfg := config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, concurrencyLimit),
		wg:                 &sync.WaitGroup{},
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	for key, value := range cfg.pages {
		fmt.Printf("%s: %d\n", key, value)
	}
}
