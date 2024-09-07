package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
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

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()
	defer func() { <-cfg.concurrencyControl }()

	cfg.concurrencyControl <- struct{}{}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Can't parse URL '%s', skipping\n", rawCurrentURL)
		return
	}
	// Avoid crawling all over the internet
	if cfg.baseURL.Hostname() != currentURL.Hostname() {
		return
	}
	key, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Can't normalize URL '%s', skipping\n", rawCurrentURL)
		return
	}
	isFirst := cfg.addPageVisit(key)
	if !isFirst {
		fmt.Printf("Already seen URL %s, skipping\n", rawCurrentURL)
		return
	}
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to parse HTML from %s: '%v'\n", rawCurrentURL, err)
		return
	}
	fmt.Printf("HTML from %s:\n%s\n", rawCurrentURL, html)
	urls, err := getURLsFromHTML(html, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Unable to extract URLs from %s: %v\n", rawCurrentURL, err)
	}
	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
