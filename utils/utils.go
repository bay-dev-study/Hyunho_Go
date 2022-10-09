package utils

import (
	"fmt"
	"log"
	"net/http"
)

func ExitIfError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ExitIfStatusCodeError(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalln(fmt.Sprintf("StatusCode : %d", resp.StatusCode))
	}
}
