package main

import (
	"fmt"
	"sort"
)

type page struct {
	internalLinks int
	url           string
}

type ByLinksThenURL []page

func (a ByLinksThenURL) Len() int { return len(a) }
func (a ByLinksThenURL) Less(i, j int) bool {
	if a[i].internalLinks == a[j].internalLinks {
		return a[i].url < a[j].url
	}
	return a[i].internalLinks < a[j].internalLinks
}
func (a ByLinksThenURL) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func printReport(pages map[string]int, baseURL string) {
	fancyPages := []page{}
	for key, val := range pages {
		fancyPages = append(fancyPages, page{
			internalLinks: val,
			url:           key,
		})
	}
	sort.Sort(sort.Reverse(ByLinksThenURL(fancyPages)))
	fmt.Printf("=============================\n  REPORT for %s\n=============================\n", baseURL)
	for _, page := range fancyPages {
		fmt.Printf("Found %d internal links to %s\n", page.internalLinks, page.url)
	}
}
