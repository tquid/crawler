package main

import (
	"fmt"
	"net/url"
)

func normalizePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	c := []rune(path)
	if c[len(c)-1] == '/' {
		return string(c[:len(c)-1]), nil
	}
	return path, nil
}

func normalizeURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("could not parse url: %v", urlStr)
	}
	normalizedPath, err := normalizePath(u.Path)
	if err != nil {
		return "", fmt.Errorf("could not normalize path '%v'", u.Path)
	}
	normalized := u.Host + normalizedPath
	return normalized, nil
}
