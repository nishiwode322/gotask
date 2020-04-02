package util

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func GetPageContext(url string) string {
	doc, _ := goquery.NewDocument(url)
	doc.Find("table").Eq(4).Find("td").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() != "" && i < 34 {
			result, _ := DecodeToGBK(selection.Text())
			fmt.Printf("++++++%d++++++\n", i)
			fmt.Println(result)
		}
	})
	return ""
}
