package main

import (
	"fmt"

	"./util"
)

func main() {

	/* 1. http echo svr
		fmt.Println("start server")
		util.EchoServer("localhost:9000")
	*/

	//2. parse province and city 
	provincecitymap, err := util.ParseProvinceAndCity("http://www.hotelaah.com/dijishi.html")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("——————————————")
		fmt.Println(provincecitymap)
	}

}
