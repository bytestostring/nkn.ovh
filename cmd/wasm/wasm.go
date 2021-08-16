package main

import (
	"fmt"
	"nknovh-wasm"
)

func main() {
	fmt.Println("Go Wasm loaded")
	c := new(nknovh_wasm.CLIENT)
	c.RegisterJSFuncs()
	c.Init()
	c.Run()
	<-make(chan bool)
}
