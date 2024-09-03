package main

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func isPathRelative(urlStr string) (bool, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false, fmt.Errorf("can't parse url '%v'", urlStr)
	}
	if parsedURL.Scheme == "" && parsedURL.Host == "" {
		return true, nil
	}
	return false, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	urls := []string{}
	z := html.NewTokenizer(strings.NewReader(htmlBody))
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return urls, nil
			}
			return nil, z.Err()
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if r, _ := isPathRelative(a.Val); r {
							urls = append(urls, rawBaseURL+a.Val)
						} else {
							urls = append(urls, a.Val)
						}
						break
					}
				}
			}
		}
	}
}
