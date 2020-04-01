package main

import (
	"fmt"
	"log"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func decodeToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(dst[:nDst]), nil
}

func getPageContext(url string) string {
	doc, _ := goquery.NewDocument(url)
	doc.Find("table").Eq(10).Find("td").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() != "" {
			result, _ := decodeToGBK(selection.Text())
			fmt.Printf("++++++%d++++++\n", i)
			fmt.Println(result)
		}
	})
	return ""
}

func main() {
	
	fmt.Println(getPageContext("http://www.hotelaah.com/dijishi.html"))
}
