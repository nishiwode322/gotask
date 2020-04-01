package main

import (
	"fmt"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var wg sync.WaitGroup
var lock sync.Mutex
var cache = struct {
	sync.Mutex
	mapping map[string]string
}{
	mapping: make(map[string]string),
}

//some data type
type point struct {
	x, y int16
}
type d3point struct {
	point
	color int
}

//some interface define
type newstructinterface interface {
	getResult() int16
}

func (d *d3point) getPointX() int16 {
	return d.x
}
func Lookup(key string) string {
	cache.Lock()
	v := cache.mapping[key]
	cache.Unlock()
	return v
}

type newint int

func getResult() (z int) {
	z = 10
	al := newint(12)
	fmt.Println("this is getResult function", al)
	return
}

func sum(values ...int) (result int) {
	for _, x := range values {
		result += x
	}
	return
}

func main() {
	a := []int{10, 20, 30}
	for i, x := range a {
		fmt.Println(i, x)
	}

	fmt.Println(getResult())
	fmt.Println(sum(3, 2, 5, 10))
	defer func() {
		fmt.Println("defer function")
	}()
	fmt.Println("hello,world")
	fmt.Println("++++++++++++++++++++++++++")
	var InterfaceInstance newstructinterface
	fmt.Println("%T",InterfaceInstance)
}
