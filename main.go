package main

import (
	"Hyunho_Go/packageScrapper"
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

	packageDataSlice := packageScrapper.GetPackageData(pageURL)
	packageScrapper.SavePackageData(packageDataSlice, "result.csv")
}
