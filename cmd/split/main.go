package main

import (
	"os"

	eso "github.com/Eurymachos/esoelog"
)

func main() {
	c := make(chan *eso.LogLine)
	go eso.LogReader(os.Args[1], c)
	eso.LogSplitter(c)
}
