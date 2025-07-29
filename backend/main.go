package main

import (
	"github.com/hasan-kilici/chat/cmd/gateaway"
	"github.com/hasan-kilici/chat/cmd/service"
)

func main() {
	go service.Start()
	gateaway.Start()
}
