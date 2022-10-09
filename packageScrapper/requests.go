package packageScrapper

import (
	"Hyunho_Go/utils"
	"net/http"
)

func GetWithErrorCheck(url string) *http.Response {
	resp, err := http.Get(url)
	utils.ExitIfError(err)
	utils.ExitIfStatusCodeError(resp)
	return resp
}
