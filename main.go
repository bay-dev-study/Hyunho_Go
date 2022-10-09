package main

import (
	"Hyunho_Go/urlcheck"
	"fmt"
)

func main() {
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
	c := make(chan urlcheck.UrlStatus)
	for _, url := range urls {
		go urlcheck.HitURL(url, c)
	}
	for range urls {
		result := <-c
		fmt.Println(result)
	}
}
