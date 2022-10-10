package scrapper

import (
	"Hyunho_Go/utils"
	"net/http"
)

func GetRequestWithErrorCheck(url string) *http.Response {
	resp, err := http.Get(url)
	utils.ExitIfError(err)
	utils.ExitIfStatusCodeError(resp)
	return resp
}
