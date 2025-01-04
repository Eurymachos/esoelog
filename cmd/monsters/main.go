package main

import (
	"flag"
	"fmt"

	eso "github.com/Eurymachos/esoelog"
)

type Ability struct {
	Unit    *eso.UnitInfo
	Ability string
}

var abilitiesUsed = make(map[Ability]bool)

var zone string

func init() {
	flag.StringVar(&zone, "z", "", "A zone name")
}

func main() {
	flag.Parse()

	for i, fn := range flag.Args() {
		log := make(chan *eso.LogLine)
		filtered := make(chan *eso.LogLine)
		events := make(chan eso.GameEvent)

		go eso.LogReader(fn, log)
		if zone == "" {
			//eso.MonsterMash()
			go eso.RunGame(log, events)
		} else {
			go eso.ZoneFilter(zone, log, filtered)
			go eso.RunGame(filtered, events)
		}
		if i > 0 {
			fmt.Println()
		}
		PrintMonsters(events)
	}
}

func PrintMonsters(msgs <-chan eso.GameEvent) {
	for msg := range msgs {
		switch m := msg.(type) {
		case *eso.EventZoneEntered:
			fmt.Println(m.Zone())

		case *eso.EventCombat:
			src := m.Source()
			if m.Ability != "" && src != nil {
				k := Ability{src.UnitInfo, m.Ability}
				abilitiesUsed[k] = true
			}

		case *eso.EventEffect:
			src := m.Source()
			if m.Ability != "" && src != nil {
				k := Ability{src.UnitInfo, m.Ability}
				abilitiesUsed[k] = true
			}

		case *eso.EventUnitSeen:
			if m.UnitState.UnitInfo.Reaction == "HOSTILE" {
				fmt.Println(m.UnitInfo)
				for k, _ := range abilitiesUsed {
					if m.UnitState.UnitInfo == k.Unit {
						fmt.Println("\t", k.Ability)
					}
				}
			}

		default:
			//fmt.Printf("%#v\n", m)
		}
	}
}
