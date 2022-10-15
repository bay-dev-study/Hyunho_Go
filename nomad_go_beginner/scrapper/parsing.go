package scrapper

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParsePackageTitle(section *goquery.Selection) (string, bool) {
	packageTitle := ""
	section.Find("a").Each(func(i int, s *goquery.Selection) {
		attrValue, exists := s.Attr("data-test-id")
		if exists && attrValue == "snippet-title" {
			title := s.Text()
			title = strings.TrimSpace(title)
			packageTitle = strings.Join(strings.Fields(title), " ")
		}
	})
	return packageTitle, (len(packageTitle) > 0)
}

func ParsePackageLink(section *goquery.Selection) (string, bool) {
	link, exists := section.Find("A").First().Attr("href")
	return link, exists
}

func ParseimportedCount(section *goquery.Selection) (int, bool) {
	resultText := ""
	section.Find("div.SearchSnippet-infoLabel").Find("a").Each(func(i int, s *goquery.Selection) {
		attr, exists := s.Attr("aria-label")
		if exists && attr == "Go to Imported By" {
			resultText = s.Find("strong").Text()
			resultText = strings.ReplaceAll(resultText, ",", "")
		}
	})

	number, err := strconv.Atoi(resultText)
	if err != nil {
		return -1, false
	}
	return number, true
}

func ParseDiscription(section *goquery.Selection) (string, bool) {
	discription := section.Find("p.SearchSnippet-synopsis").Text()
	discription = strings.TrimSpace(discription)
	return discription, (len(discription) > 0)
}

func ParseFuncInfo(section *goquery.Selection) (string, bool) {
	funcInfo := section.Find("div.Documentation-declaration").Text()
	funcInfo = strings.TrimSpace(funcInfo)
	return funcInfo, (len(funcInfo) > 0)
}
