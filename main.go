package main

import (
	"fmt"
	"os"
)

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
	pages := map[string]int{}
	crawlPage(rawBaseURL, rawBaseURL, pages)
	for key, value := range pages {
		fmt.Printf("%s: %d\n", key, value)
	}
}
