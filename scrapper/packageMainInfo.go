package scrapper

import (
	"Hyunho_Go/utils"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type PackageMainInfo struct {
	packageName   string
	url           string
	discription   string
	importedCount int
}

func (packageMainInfo *PackageMainInfo) Url() string {
	return packageMainInfo.url
}

func (packageMainInfo *PackageMainInfo) GetCSVHead() []string {
	return []string{"packageName", "url", "discription", "importedCount"}
}

func (packageMainInfo *PackageMainInfo) ToCSVForm() []string {
	return []string{packageMainInfo.packageName, packageMainInfo.url, packageMainInfo.discription, strconv.Itoa(packageMainInfo.importedCount)}
}

func GetPackageMainInfo(pageURL string, returnPackageMainInfoChannel chan<- []PackageMainInfo) {
	var packageMainInfoSlice []PackageMainInfo

	response := GetRequestWithErrorCheck(pageURL)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	utils.ExitIfError(err)

	doc.Find("div.SearchSnippet").Each(func(i int, s *goquery.Selection) {
		packageName, existsName := ParsePackageTitle(s)

		url, existsLink := ParsePackageLink(s)

		discription, existsDiscription := ParseDiscription(s)

		importedCount, existsImportedCount := ParseimportedCount(s)

		if existsName && existsLink && existsDiscription && existsImportedCount {
			packageMainInfoSlice = append(packageMainInfoSlice, PackageMainInfo{packageName: packageName, url: url, discription: discription, importedCount: importedCount})
		}
	})
	returnPackageMainInfoChannel <- packageMainInfoSlice
}
