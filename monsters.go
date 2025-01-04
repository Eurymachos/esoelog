package esoelog

import (
	"fmt"
)

type monster struct {
	name      string
	healthMax int
}

var seen map[monster]bool
var zone string

func MonsterMash(c <-chan *LogLine) {

	// A map of unitId to name
	monsters := make(map[int]*UnitInfo)

	seen = make(map[monster]bool)

	for l := range c {
		switch l.LineType {
		case ZoneChanged:
			// Unit IDs are reset across zone changes
			monsters = make(map[int]*UnitInfo)
			if l.LineData[4] == "NONE" {
				zone = fmt.Sprintf("%s", l.LineData[3])
			} else {
				zone = fmt.Sprintf("%s (%s)", l.LineData[3], l.LineData[4])
			}

		case UnitAdded:

			// 1: __UNIT_ADDED__ - unitId, unitType,
			// 4: isLocalPlayer, playerPerSessionId, monsterId,
			// 7: isBoss, classId, raceId, name, displayName,
			// 12: characterId, level, championPoints, ownerUnitId,
			// 16: reaction, isGroupedWithLocalPlayer

			ui := NewUnitInfo(l.LineData[2:])
			if ui.Reaction != "HOSTILE" {
				continue
			}

			if check, found := monsters[ui.unitID]; found {
				fmt.Println("Found", ui.unitID, ui.name, "when I already had", check.unitID, check.name)
			}

			monsters[ui.unitID] = ui

		case BeginCast, EffectChanged:
			checkUnitState(monsters, l.LineData[6:])
			if l.LineData[16] != "*" {
				checkUnitState(monsters, l.LineData[16:])
			}

		case CombatEvent:
			checkUnitState(monsters, l.LineData[9:])
			if l.LineData[19] != "*" {
				checkUnitState(monsters, l.LineData[16:])
			}
		}
	}
}

func checkUnitState(monsters map[int]*UnitInfo, data []string) {
	us := parseUnitState(data)

	if ui, found := monsters[us.unitID]; found {
		// We only need to see the max health once
		delete(monsters, us.unitID)

		if !seen[monster{ui.name, us.healthMax}] {
			if ui.isBoss {
				fmt.Println(zone, ui.name, "(Boss)", us.healthMax)
			} else {
				fmt.Println(zone, ui.name, us.healthMax)
			}
			seen[monster{ui.name, us.healthMax}] = true
		}
	}
}
