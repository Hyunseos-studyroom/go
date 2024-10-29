package main

import (
	"fmt"
	"net/http"
)

func main() {
	var results = make(map[string]string)

	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://academy.nomadcoders.co/",
	}

	for _, url := range urls {
		result := "OK"
		err := higURL(url)
		if err != nil {
			result = "ERROR"
		}
		results[url] = result
	}
	for url, result := range results {
		fmt.Println(url, result)
	}
}

func higURL(url string) (err error) {
	fmt.Println("CHecking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return err
	}
	return nil
}
