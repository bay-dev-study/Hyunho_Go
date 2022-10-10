package main

import (
	"Hyunho_Go/scrapper"
	"Hyunho_Go/utils"

	"github.com/labstack/echo/v4"
)

const zipFileName = "result.zip"

func handleScrape(c echo.Context) error {
	searchWord := c.FormValue("packageName")
	fileNames := scrapper.ScrapePackageData(searchWord)
	utils.CompressFiles(fileNames, zipFileName)
	defer utils.RemoveFiles(append(fileNames, zipFileName))
	return c.Attachment(zipFileName, zipFileName)
}

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func main() {

	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
