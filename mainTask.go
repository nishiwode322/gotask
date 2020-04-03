package main

import (
	"fmt"

	"./util"
)

func main() {
	provincecitymap, err := util.ParseProvinceAndCity("http://www.hotelaah.com/dijishi.html")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("——————————————")
		fmt.Println(provincecitymap)
	}
}
