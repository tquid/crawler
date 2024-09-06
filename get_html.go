package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawUrl string) (string, error) {
	resp, err := http.Get(rawUrl)
	if err != nil {
		return "", fmt.Errorf("error getting url '%s': %v", rawUrl, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error from url '%s': %s", rawUrl, resp.Status)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("content-type is not text/html: got %s", contentType)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	return string(content), nil
}
