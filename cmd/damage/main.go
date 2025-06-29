package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"unicode/utf8"

	eso "github.com/Eurymachos/esoelog"
)

// |DateTimestamp|Source (player/pet/companion)|Skill/Attack used|Target|Damage|(A 'Crit?' Boolean would be handy but not essential).

var zone string
var delimiter string

func init() {
	flag.StringVar(&zone, "z", "", "A zone name")
	flag.StringVar(&delimiter, "d", "", "Delimiter character")
}

func main() {
	flag.Parse()

	for i, fn := range flag.Args() {
		if i > 0 {
			fmt.Println()
		}

		log := make(chan *eso.LogLine)
		filtered := make(chan *eso.LogLine)
		events := make(chan eso.GameEvent)

		go eso.LogReader(fn, log)
		go eso.ZoneFilter(zone, log, filtered)
		go eso.RunGame(filtered, events)
		PrintFight(events)
	}
}

func PrintFight(msgs <-chan eso.GameEvent) {

	out := csv.NewWriter(os.Stdout)

	if delimiter != "" {
		out.Comma, _ = utf8.DecodeRuneInString(delimiter)

	}

	for msg := range msgs {
		rec := make([]string, 7)

		switch m := msg.(type) {
		case *eso.EventCombat:
			switch m.ActionResult {
			case "DAMAGE":
				rec[3] = "DAMAGE"
				rec[6] = "F"
			case "CRITICAL_DAMAGE":
				rec[3] = "DAMAGE"
				rec[6] = "T"
			case "DOT_TICK":
				rec[3] = "DOT_TICK"
				rec[6] = "F"
			case "DOT_TICK_CRITICAL":
				rec[3] = "DOT_TICK"
				rec[6] = "T"
			default:
				continue
			}

			src := m.Source()
			dst := m.Target()

			if src.Reaction != "PLAYER_ALLY" &&
				src.Reaction != "COMPANION" &&
				src.Reaction != "NPC_ALLY" {
				continue
			}

			if src.Reaction == "PLAYER_ALLY" && src.ID() != 1 {
				continue
			}

			rec[0] = m.When
			if src != nil {
				rec[1] = src.Name()
			}
			rec[2] = m.Ability
			if dst != nil {
				rec[4] = dst.Name()
			}
			rec[5] = m.HitValue

			err := out.Write(rec)
			if err != nil {
				fmt.Println("CSV error:", err)
				fmt.Println("Couldn't print:", rec)
			}
		}
	}

	out.Flush()
}
