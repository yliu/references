package main

import (
	"fmt"
	"os"
)

var isCtags bool
var isGlobal bool
var extendTag bool

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

	extendTag = true
}
func main() {
	newUI(os.Args[1:])
}
