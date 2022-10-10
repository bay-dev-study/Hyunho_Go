package scrapper

import "fmt"

const mainURL = "https://pkg.go.dev"
const mainInfoFileName = "mainInfo.csv"
const funcInfoFileName = "funcInfo.csv"

func ScrapePackageData(searchWord string) []string {
	pageURL := fmt.Sprintf("%s/search?q=%s", mainURL, searchWord)
	fmt.Println(pageURL)

	returnPackageMainInfoChannel := make(chan []PackageMainInfo)
	go GetPackageMainInfo(pageURL, returnPackageMainInfoChannel)
	packageMainInfoSlice := <-returnPackageMainInfoChannel
	SavePackageMainInfo(packageMainInfoSlice, mainInfoFileName)

	returnPackageFuncInfoChannel := make(chan PackageFuncInfo)
	for _, packageMainInfo := range packageMainInfoSlice {
		go GetPackageFuncInfo(mainURL+packageMainInfo.Url(), packageMainInfo, returnPackageFuncInfoChannel) // packageMainInfo를 pointer로 넘겨주면 왜 안 되는것일까?
	}
	var packageFuncInfoSlice []PackageFuncInfo
	for range packageMainInfoSlice {
		packageFuncInfoSlice = append(packageFuncInfoSlice, <-returnPackageFuncInfoChannel)
	}
	SavePackageFuncInfo(packageFuncInfoSlice, funcInfoFileName)
	return []string{mainInfoFileName, funcInfoFileName}
}
