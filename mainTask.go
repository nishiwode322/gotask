package main

import (
	"./util"
)

func main() {

	// 1. http echo svr
	/*
		fmt.Println("start server")
		util.EchoServer("localhost:9000")
	*/

	//2. parse province and city from given web page
	/*
		provinceCityMap, err := util.ParseProvinceAndCity("http://www.hotelaah.com/dijishi.html")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("—————————————")
			fmt.Println(provinceCityMap)
		}
	*/

	//3.just some function test
	//3.1 test go connect kafka
	/*
		util.RunKafkaSample()
	*/

	//3.2 test go connect mysql
	util.RunMysqlSample()

	//3.3 TODO: go connect redis
}
