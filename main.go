package main

import (
	"fmt"
	"os"
)

var isCtags bool
var isGlobal bool

func init() {
	e := checkCtags()
	if e == nil {
		isCtags = true
	}
	e = checkGlobal()
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(0)
	}
	isGlobal = true
}
func main() {
	newUI()
}
