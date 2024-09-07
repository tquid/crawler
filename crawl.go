package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	baseURL            *url.URL
	concurrencyControl chan struct{}
	maxPages           int
	mu                 *sync.Mutex
	pages              map[string]int
	wg                 *sync.WaitGroup
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	_, exists := cfg.pages[normalizedURL]
	if !exists {
		cfg.pages[normalizedURL] = 1
	} else {
		cfg.pages[normalizedURL] += 1
	}
	cfg.mu.Unlock()
	return !exists
}

func (cfg *config) maxPagesReached() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages) >= cfg.maxPages
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()
	defer func() { <-cfg.concurrencyControl }()

	cfg.concurrencyControl <- struct{}{}

	if cfg.maxPagesReached() {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Can't parse URL '%s', skipping\n", rawCurrentURL)
		return
	}
	// Avoid crawling all over the internet
	if cfg.baseURL.Hostname() != currentURL.Hostname() {
		return
	}
	_, err = normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Can't normalize URL '%s', skipping\n", rawCurrentURL)
		return
	}
	isFirst := cfg.addPageVisit(rawCurrentURL)
	if !isFirst {
		return
	}
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to parse HTML from %s: '%v'\n", rawCurrentURL, err)
		return
	}
	urls, err := getURLsFromHTML(html, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Unable to extract URLs from %s: %v\n", rawCurrentURL, err)
	}
	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
