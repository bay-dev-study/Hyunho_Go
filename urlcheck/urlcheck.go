package urlcheck

import "net/http"

type UrlStatus struct {
	url    string
	status string
}

func HitURL(url string, c chan<- UrlStatus) {
	response, err := http.Get(url)
	if err != nil || response.StatusCode >= 400 {
		c <- UrlStatus{url: url, status: "Failed"}
	} else {
		c <- UrlStatus{url: url, status: "SUCCEED"}
	}
}
