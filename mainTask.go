package main

import (
	"fmt"
	"./util"
)


func main() {
	//var provinceAndCityMap = make(map[string][]string)
	//provinceAndCityMap["江苏"]=[]string{"南京","丹阳"}
	//fmt.Println(provinceAndCityMap["江苏"])
	//provinceAndCityMap["江苏"]=append(provinceAndCityMap["江苏"],"苏州")
	//fmt.Println(provinceAndCityMap["江苏"])

	fmt.Println(util.GetPageContext("http://www.hotelaah.com/dijishi.html"))
}
