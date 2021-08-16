package main

import (
	"log"
	x "nknovh-engine"
)

func main() {
	var e x.NKNOVH
		if err := e.Build(); err != nil {
			log.Fatal("Build() returned error: ", err.Error())
		}
	go func() {
		if err := e.Run(); err != nil {
			log.Fatal("Run() returned error: ", err.Error())
		}
	}()

	e.Listen()

}
