package util

import (
	"errors"
	"fmt"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var BaseURL string = "http://www.hotelaah.com/"

func ParseProvinceAndCity(url string) (map[string][]string, error) {
	result := make(map[string][]string)
	wg := sync.WaitGroup{}
	wg.Add(27)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("url is invalid:" + url)
	}
	//index is 1,2,24,26,31,32,33
	doc.Find("table").Eq(4).Find("a").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() != "" && i < 34 {
			provinceName, _ := DecodeGBK(selection.Text())
			fmt.Println(provinceName)
			if i == 1 || i == 2 || i == 24 || i == 26 || i == 31 || i == 32 || i == 33 {
				result[provinceName] = []string{provinceName}
			} else {
				//get suburl
				tempURL, _ := selection.Attr("href")
				subURL := BaseURL + tempURL
				result[provinceName] = []string{}
				go parseCity(subURL, provinceName, &result, &wg)
			}
		}
	})
	wg.Wait()
	return result, nil
}

func parseCity(url string, provinceName string, provinceCity *map[string][]string, waitGroup *sync.WaitGroup) {
	result := []string{}
	doc, _ := goquery.NewDocument(url)
	doc.Find("table").Eq(10).Find("td").Each(func(i int, selection *goquery.Selection) {
		if i == 3 || i == 4 {
			targetSelection := selection.Find("a")
			selectionSize := targetSelection.Size()
			targetSelection.Each(func(i int, subSelection *goquery.Selection) {
				if subSelection.Text() != "" && i < selectionSize-1 {
					cityName, _ := DecodeGBK(subSelection.Text())
					result = append(result, cityName)
				}
			})
		}
	})
	(*provinceCity)[provinceName] = result
	waitGroup.Done()
}
