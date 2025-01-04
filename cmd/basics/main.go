// Basics verifies certain assumptions about how ESO encounter logs work.
// It does this by calculating some basic statistics and printing them out so
// that the user can manually verify those assumptions.
//
// These assumptions are:
//
//  1. The log file is always closed with an END_LOG line.
//  2. Unit IDs (and therefore units) do not persist across a zone change.
//  3. Unit IDs (and therefore units) do persist across a map change.
//  4. Ability definitions persist across zone changes.
//
// Running this program over an encounter log will output something like this:
//
//	$ go run ./cmd/basics/main.go ~/Encounter-2024-11-22.log
//	Abilities used undefined in zone 133107
//	Abilities used undefined in log 0
//	Map changes with 0 extant units 0
//	Map changes with >0 extant units 110
//	Zone changes with 0 extant units 51
//	Zone changes with >0 extant units 0
//	Log ends with 0 extant units 14
//	Log ends with >0 extant units 0
//	Log trailing lines 0
//
// The encounter log has UNIT_ADDED and UNIT_REMOVED lines and any
// unit that has been added but not removed is "extant". I can see
// that many units are extant during a map change, but never during a
// zone change or when a log ends. I infer from that, that units
// persist across map changes, but not zone or log changes.
//
// The program also counts how many lines follow the last END_LOG just
// in case a log file hasn't been closed off properly.
//
// Finally the use of abilities is counted, tracking how many are used
// without an ABILITY_INFO record in that zone.
//
// You can run this program with a -v flag to output some of the data
// being collected in order to try and get some understanding of the
// underlying processing.
package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	eso "github.com/Eurymachos/esoelog"
)

const (
	ZoneUndefAbilities = iota
	LogUndefAbilities
	MapChgZeroUnits
	MapChgMoreUnits
	ZoneChgZeroUnits
	ZoneChgMoreUnits
	LogEndZeroUnits
	LogEndMoreUnits
	TrailingLines
	NumStats
)

var msgs = []string{
	"Abilities used undefined in zone",
	"Abilities used undefined in log",
	"Map changes with 0 extant units",
	"Map changes with >0 extant units",
	"Zone changes with 0 extant units",
	"Zone changes with >0 extant units",
	"Log ends with 0 extant units",
	"Log ends with >0 extant units",
	"Log trailing lines",
}

var statsFlag, verboseFlag bool

func init() {
	flag.BoolVar(&statsFlag, "s", true, "output summary statistics")
	flag.BoolVar(&verboseFlag, "v", false, "output individual counts ")
}

func main() {
	flag.Parse()

	for i, fn := range flag.Args() {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println("Log file:", fn)
		c := make(chan *eso.LogLine)
		go eso.LogReader(fn, c)
		DoIt(c)
	}
}

func DoIt(c <-chan *eso.LogLine) {
	stats := make(map[int]int)
	zoneA := make(map[string]bool)     // Zone defined abilities
	logA := make(map[string]bool)      // Log defined abilities
	logUndefA := make(map[string]bool) // Log undefined abilities

	lineCount := 0
	units := make(map[string][]string)

	for l := range c {
		lineCount++
		switch l.LineType {
		case eso.BeginLog:
			if verboseFlag {
				epochMs, err := strconv.ParseInt(l.LineData[2], 10, 64)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(l, time.UnixMilli(epochMs))
			}

			logA = make(map[string]bool)
			zoneA = make(map[string]bool)

		case eso.EndLog:
			if len(units) == 0 {
				stats[LogEndZeroUnits]++
			} else {
				stats[LogEndMoreUnits]++
			}

			if verboseFlag {
				fmt.Println(l, len(units), "units")
				fmt.Println(l, lineCount, "lines")
			}

			lineCount = 0
			units = make(map[string][]string)

		case eso.UnitAdded:
			units[l.LineData[2]] = l.LineData

		case eso.UnitRemoved:
			delete(units, l.LineData[2])

		case eso.MapChanged:
			if len(units) == 0 {
				stats[MapChgZeroUnits]++
			} else {
				stats[MapChgMoreUnits]++
			}

		case eso.ZoneChanged:
			if len(units) == 0 {
				stats[ZoneChgZeroUnits]++
			} else {
				stats[ZoneChgMoreUnits]++
			}

			if verboseFlag {
				fmt.Println(l, len(units), "units")
			}

			zoneA = make(map[string]bool)

		case eso.BeginCast, eso.CombatEvent, eso.EffectChanged, eso.EffectInfo:
			var abilityID string
			switch l.LineType {
			case eso.BeginCast:
				abilityID = l.LineData[5]
			case eso.CombatEvent:
				abilityID = l.LineData[8]
			case eso.EffectChanged:
				abilityID = l.LineData[5]
			case eso.EffectInfo:
				abilityID = l.LineData[2]
			}
			// Soul gem resurrections have no ability ID
			if abilityID == "0" {
				continue
			}
			if !zoneA[abilityID] {
				stats[ZoneUndefAbilities]++
			}
			if !logA[abilityID] {
				stats[LogUndefAbilities]++
				logUndefA[abilityID] = true
			}

		case eso.AbilityInfo:
			zoneA[l.LineData[2]] = true
			logA[l.LineData[2]] = true
		}
	}

	if verboseFlag {
		fmt.Println(lineCount, "trailing lines")
	}

	if statsFlag && verboseFlag {
		fmt.Println()
	}

	if statsFlag {
		for i := range NumStats {
			fmt.Println(msgs[i], stats[i])
		}

		if len(logUndefA) > 0 {
			fmt.Println("Log undefined ability IDs:", logUndefA)
		}
	}
}
