package util

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func ParseProvinceAndCity(url string) (map[string]CityList, error) {
	result := make(map[string]CityList)
	//get base url
	index := strings.LastIndex(url, "/")
	baseURL := url[:index+1]
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
				result[provinceName] = CityList{provinceName}
			} else {
				//get suburl
				tempURL, _ := selection.Attr("href")
				subURL := baseURL + tempURL
				result[provinceName] = CityList{}
				go parseCity(subURL, provinceName, result, &wg)
			}
		}
	})
	wg.Wait()
	return result, nil
}

func parseCity(url string, provinceName string, provinceCity map[string]CityList, waitGroup *sync.WaitGroup) {
	result := CityList{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return
	}
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
	provinceCity[provinceName] = result

	defer waitGroup.Done()
}
