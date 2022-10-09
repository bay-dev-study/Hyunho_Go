package packageScrapper

import (
	"Hyunho_Go/utils"
	"bufio"
	"encoding/csv"
	"os"
)

func SavePackageData(packageDataSlice []PackageData, fileName string) {
	file, err := os.Create("./" + fileName)
	utils.ExitIfError(err)
	defer file.Close()

	wr := csv.NewWriter(bufio.NewWriter(file))

	var records [][]string

	records = append(records, packageDataSlice[0].GetCSVHead())
	for _, packageData := range packageDataSlice {
		records = append(records, packageData.ToCSVForm())
	}
	wr.WriteAll(records) // calls Flush internally
}
