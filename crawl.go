package main

import (
	"fmt"
	"net/url"
	"os"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Can't parse base URL '%s', please check it and try again\n", rawBaseURL)
		os.Exit(1)
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Can't parse URL '%s', skipping\n", rawCurrentURL)
		return
	}
	// Avoid crawling all over the internet
	if baseURL.Hostname() != currentURL.Hostname() {
		return
	}
	key, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Can't normalize URL '%s', skipping\n", rawCurrentURL)
		return
	}
	value, exists := pages[key]
	if exists {
		pages[key] = value + 1
		fmt.Printf("Already seen URL %s, skipping\n", rawCurrentURL)
		return
	} else {
		pages[key] = 1
	}
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to parse HTML from %s: '%v'\n", rawCurrentURL, err)
		return
	}
	fmt.Printf("HTML from %s:\n%s\n", rawCurrentURL, html)
	urls, err := getURLsFromHTML(html, rawBaseURL)
	if err != nil {
		fmt.Printf("Unable to extract URLs from %s: %v\n", rawCurrentURL, err)
	}
	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}
