package main

import (
	"Hyunho_Go/scrapper"
	"Hyunho_Go/utils"
	"fmt"
)

func inputSearchWord() (searchWord string, err error) {
	fmt.Printf("Search Word: ")
	_, err = fmt.Scanf("%s", &searchWord)
	return
}

func main() {
	mainURL := "https://pkg.go.dev"
	searchWord, err := inputSearchWord()
	utils.ExitIfError(err)

	pageURL := fmt.Sprintf("%s/search?q=%s", mainURL, searchWord)
	fmt.Println(pageURL)

	returnPackageMainInfoChannel := make(chan []scrapper.PackageMainInfo)
	go scrapper.GetPackageMainInfo(pageURL, returnPackageMainInfoChannel)
	packageMainInfoSlice := <-returnPackageMainInfoChannel
	scrapper.SavePackageMainInfo(packageMainInfoSlice, "mainInfo.csv")

	returnPackageFuncInfoChannel := make(chan scrapper.PackageFuncInfo)
	for _, packageMainInfo := range packageMainInfoSlice {
		go scrapper.GetPackageFuncInfo(mainURL+packageMainInfo.Url(), packageMainInfo, returnPackageFuncInfoChannel) // packageMainInfo를 pointer로 넘겨주면 왜 안 되는것일까?
	}
	var packageFuncInfoSlice []scrapper.PackageFuncInfo
	for range packageMainInfoSlice {
		packageFuncInfoSlice = append(packageFuncInfoSlice, <-returnPackageFuncInfoChannel)
	}
	scrapper.SavePackageFuncInfo(packageFuncInfoSlice, "funcInfo.csv")
}
