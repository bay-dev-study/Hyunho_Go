package packageScrapper

import (
	"Hyunho_Go/utils"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type PackageData struct {
	packageName   string
	url           string
	discription   string
	importedCount int
}

func (packageData *PackageData) GetCSVHead() []string {
	return []string{"packageName", "url", "discription", "importedCount"}
}

func (packageData *PackageData) ToCSVForm() []string {
	return []string{packageData.packageName, packageData.url, packageData.discription, strconv.Itoa(packageData.importedCount)}
}

func GetPackageData(pageURL string) []PackageData {
	var packageDataSlice []PackageData

	response := GetWithErrorCheck(pageURL)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	utils.ExitIfError(err)

	doc.Find("div.SearchSnippet").Each(func(i int, s *goquery.Selection) {
		packageName, existsName := ParsePackageTitle(s)

		url, existsLink := ParsePackageLink(s)

		discription, existsDiscription := ParseDiscription(s)

		importedCount, existsImportedCount := ParseimportedCount(s)

		if existsName && existsLink && existsDiscription && existsImportedCount {
			packageDataSlice = append(packageDataSlice, PackageData{packageName: packageName, url: url, discription: discription, importedCount: importedCount})
		}
	})
	return packageDataSlice
}
