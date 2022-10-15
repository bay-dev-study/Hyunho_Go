package scrapper

import (
	"Hyunho_Go/utils"
	"bufio"
	"encoding/csv"
	"os"
)

func writeDataToCSV(fileName string, records [][]string) error {
	file, err := os.Create("./" + fileName)
	defer file.Close()
	utils.ExitIfError(err)
	wr := csv.NewWriter(bufio.NewWriter(file))
	return wr.WriteAll(records)
}

func SavePackageMainInfo(packageMainInfoSlice []PackageMainInfo, fileName string) {
	var records [][]string
	records = append(records, packageMainInfoSlice[0].GetCSVHead())
	for _, packageMainInfo := range packageMainInfoSlice {
		records = append(records, packageMainInfo.ToCSVForm())
	}
	err := writeDataToCSV(fileName, records)
	utils.ExitIfError(err)
}

func SavePackageFuncInfo(packageFuncInfoSlice []PackageFuncInfo, fileName string) {
	var records [][]string
	records = append(records, packageFuncInfoSlice[0].GetCSVHead())
	for _, packageFuncInfo := range packageFuncInfoSlice {
		records = append(records, packageFuncInfo.ToCSVForm()...)
	}
	err := writeDataToCSV(fileName, records)
	utils.ExitIfError(err)
}
