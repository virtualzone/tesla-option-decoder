package main

import (
	"os"
)

func main() {
	InitOptionCodeMap()
	ServeHTTP()
	os.Exit(0)
}
