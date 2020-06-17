package main

import (
	"github.com/anorb/lucrum"
)

func main() {
	luc := lucrum.Init()

	go luc.UpdateLoop()
	luc.Run()
}
