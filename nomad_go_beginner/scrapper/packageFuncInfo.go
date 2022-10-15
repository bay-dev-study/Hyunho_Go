package scrapper

import (
	"Hyunho_Go/utils"

	"github.com/PuerkitoBio/goquery"
)

type PackageFuncInfo struct {
	packageName string
	url         string
	funcInfo    []string
}

func (packageFuncInfo *PackageFuncInfo) GetCSVHead() []string {
	return []string{"packageName", "url", "funcInfo"}
}

func (packageFuncInfo *PackageFuncInfo) ToCSVForm() [][]string {
	var records [][]string
	for _, funcInfo := range packageFuncInfo.funcInfo {
		records = append(records, []string{packageFuncInfo.packageName, packageFuncInfo.url, funcInfo})
	}
	return records
}

func GetPackageFuncInfo(pageURL string, packageMainInfo PackageMainInfo, returnPackageFuncInfoChannel chan<- PackageFuncInfo) {
	var funcInfoSlice []string
	response := GetRequestWithErrorCheck(pageURL)
	doc, err := goquery.NewDocumentFromReader(response.Body)
	utils.ExitIfError(err)

	doc.Find("div.Documentation-function").Each(func(i int, s *goquery.Selection) {
		funcInfo, exists := ParseFuncInfo(s)
		if exists {
			funcInfoSlice = append(funcInfoSlice, funcInfo)
		}
	})
	returnPackageFuncInfoChannel <- PackageFuncInfo{packageName: packageMainInfo.packageName, url: packageMainInfo.url, funcInfo: funcInfoSlice}
}
