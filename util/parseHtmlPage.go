package util

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var BaseURL string = "http://www.hotelaah.com/"

func ParseProvinceAndCity(url string) (map[string][]string, error) {
	result := make(map[string][]string)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("url is invalid:" + url)
	}
	//index is 1,2,24,26,31,32,33
	doc.Find("table").Eq(4).Find("a").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() != "" && i < 34 {
			provincename, _ := DecodeGBK(selection.Text())
			fmt.Println("省:", provincename)
			if i == 1 || i == 2 || i == 24 || i == 26 || i == 31 || i == 32 || i == 33 {
				result[provincename] = []string{provincename}
			} else {
				//get suburl
				tempurl, _ := selection.Attr("href")
				suburl := BaseURL + tempurl
				result[provincename], _ = ParseCity(suburl)
			}
		}
	})
	return result, nil
}

func ParseCity(url string) ([]string, error) {
	var result []string
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("url is invalid:" + url)
	}
	doc.Find("table").Eq(10).Find("td").Each(func(i int, selection *goquery.Selection) {
		if i == 3 || i == 4 {
			targetselection := selection.Find("a")
			selectionsize := targetselection.Size()
			targetselection.Each(func(i int, subselection *goquery.Selection) {
				if subselection.Text() != "" && i < selectionsize-1 {
					cityname, _ := DecodeGBK(subselection.Text())
					result = append(result, cityname)
					fmt.Println(cityname)
				}
			})
		}
	})
	return result, nil
}
