package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// NOTE: img tags invalid tokens but that must mean they're content tokens
// figure out the enclosing tag, should be valid, then becomes a matter of extracting
// from broken img tag token

// fetchHTML fetches the provided URL and returns the response body or an error
// func fetchHTML(URL string) (io.ReadCloser, error) {}

// extractTitle returns the content within the <title> element or an error
// func extractTitle(body io.ReadCloser) (string, error) {}

//fetchTitle fetches the page title for a UR
// func fetchTitle(URL string) (string, error) {}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage:\n pagetitle <url>\n")
		os.Exit(1)
	}

	URL := os.Args[1]

	resp, err := http.Get(URL) // resp is a pointer to the http.Response struct

	if err != nil {
		log.Fatalf("error fetching URL: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code was %d\n", resp.StatusCode)
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		log.Fatalf("response content type was %s not text/html\n", ctype)
	}

	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		}

		name, hasAttr := tokenizer.TagName()
		if hasAttr {
			s := string(name)
			if s == "img" {
				_, val, _ := tokenizer.TagAttr()
				fmt.Println("Image path is:", string(val))
			}
		}
	}
	defer resp.Body.Close()
}
