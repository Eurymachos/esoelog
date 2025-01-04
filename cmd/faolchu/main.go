package main

import (
	"flag"
	"fmt"

	eso "github.com/Eurymachos/esoelog"
)

var zone string

func init() {
	flag.StringVar(&zone, "z", "Camlorn Keep", "A zone name")
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

	for msg := range msgs {
		switch m := msg.(type) {
		case *eso.EventCombat:
			if m.ActionResult == "ABILITY_ON_COOLDOWN" ||
				m.ActionResult == "QUEUED" {
				continue
			}

			src := m.Source()
			dst := m.Target()

			fmt.Printf("%s %s %s %s %s",
				m.When, m.ActionResult,
				m.DamageType, m.HitValue, m.Ability)
			if src != nil {
				fmt.Print(" ", src.UnitInfo)
			}
			if dst != nil && dst.UnitInfo != nil {
				fmt.Println("", "->", dst)
			} else {
				fmt.Println()
			}

		case *eso.EventEffect:
			src := m.Source()
			dst := m.Target()
			fmt.Printf("%s %s %s", m.When, m.ChangeType, m.Ability)
			if src != nil {
				fmt.Print(" ", src.UnitInfo)
			}
			fmt.Println(" ->", dst)

		case *eso.EventZoneEntered:
		case *eso.EventZoneExited:
		case *eso.EventUnitSeen:

		default:
			fmt.Printf("%#v\n", m)
		}
	}
}
