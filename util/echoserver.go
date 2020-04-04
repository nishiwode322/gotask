package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func echohandler(writer http.ResponseWriter, read *http.Request) {
	clientdata, _ := ioutil.ReadAll(read.Body)
	fmt.Println("server receive data:", clientdata)
	writer.Write(clientdata)
}

func EchoServer(ipandport string) {
	http.HandleFunc("/", echohandler)
	fmt.Println("echo server is listening:" + ipandport)
	if err := http.ListenAndServe(ipandport, nil); err != nil {
		fmt.Println(err)
	}
}
