package main

import (
	"os"

	"github.com/Eurymachos/esoelog"
)

func main() {
	lines := make(chan *esoelog.LogLine)
	events := make(chan esoelog.GameEvent)

	go esoelog.LogReader(os.Args[1], lines)
	go esoelog.RunGame(lines, events)
	esoelog.FindBoundary(events)
}
