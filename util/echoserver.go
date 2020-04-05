package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func echoHandler(writer http.ResponseWriter, read *http.Request) {
	clientData, _ := ioutil.ReadAll(read.Body)
	fmt.Println("server receive data:", clientData)
	writer.Write(clientData)
}

func EchoServer(ipAndPort string) {
	http.HandleFunc("/", echoHandler)
	fmt.Println("echo server is listening:" + ipAndPort)
	if err := http.ListenAndServe(ipAndPort, nil); err != nil {
		fmt.Println("ListenAndServe err! ", err)
	}
}
