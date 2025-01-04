package esoelog

import (
	"fmt"
	"strconv"
)

// UNIT_ADDED, unitId, unitType, isLocalPlayer, playerPerSessionId,
// monsterId, isBoss, classId, raceId, name, displayName, characterId,
// level, championPoints, ownerUnitId, reaction,
// isGroupedWithLocalPlayer
//
// 4520901,UNIT_ADDED,5,MONSTER,F,0,15134,F,0,0,"Faolchu the
// Reborn","",0,50,160,0,HOSTILE,F

// UnitInfo holds the immutable state of a unit.
type UnitInfo struct {
	unitID                   int
	unitType                 string
	isLocalPlayer            bool
	playerPerSessionID       int
	monsterID                int
	isBoss                   bool
	classID                  int
	raceID                   int
	name                     string
	displayName              string
	characterID              int64
	level                    int
	championPoints           int
	ownerUnitID              int
	Reaction                 string
	isGroupedWithLocalPlayer bool
	maxHealth                []int // transformations make more than one
}

func NewUnitInfo(line []string) *UnitInfo {
	ui := new(UnitInfo)

	ui.unitID = mustInt(line[0])
	ui.unitType = line[1]
	ui.isLocalPlayer = line[2] == "T"
	ui.playerPerSessionID = mustInt(line[3])
	ui.monsterID = mustInt(line[4])
	ui.isBoss = line[5] == "T"
	ui.classID = mustInt(line[6])
	ui.raceID = mustInt(line[7])
	ui.name = line[8]
	ui.displayName = line[9]
	ui.characterID, _ = strconv.ParseInt(line[10], 10, 64)
	ui.level = mustInt(line[11])
	ui.championPoints = mustInt(line[12])
	ui.ownerUnitID = mustInt(line[13])
	ui.Reaction = line[14]
	ui.isGroupedWithLocalPlayer = line[15] == "T"

	//ui.name = fmt.Sprintf("%s[%d]", ui.name, ui.unitId)

	return ui
}

func (ui *UnitInfo) String() string {
	name := ui.Name()

	// Player max health can change a lot and it isn't useful on the wiki
	// so just don't report it.
	if ui.unitType == "PLAYER" {
		return name
	}

	if len(ui.maxHealth) == 0 {
		return name + " [?health]"
	}

	return fmt.Sprintf("%s %v", name, ui.maxHealth)
}

func (ui *UnitInfo) Name() string {
	if ui.name == "" {
		return ui.unitType
	}
	return ui.name
}

func (ui *UnitInfo) ID() int {
	return ui.unitID
}

func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	if strconv.Itoa(i) != s {
		panic(fmt.Sprintf("%s seems to be bigger than an int", s))
	}
	return i
}

// UnitState holds the mutable state of a unit.
type UnitState struct {
	*UnitInfo
	unitID                int
	health, healthMax     int
	magicka, magickaMax   int
	stamina, staminaMax   int
	ultimate, ultimateMax int
	werewolf, werewolfMax int
	shield                string
	mapNX, mapNY          float64
	heading               float64
}

// UnitState refers to the following fields for a unit: unitId,
// health/max, magicka/max, stamina/max, ultimate/max, werewolf/max,
// shield, map NX, map NY, headingRadians
func parseUnitState(line []string) *UnitState {
	u := new(UnitState)
	u.unitID, _ = strconv.Atoi(line[0])
	if u.unitID == 0 {
		return nil
	}
	fmt.Sscanf(line[1], "%d/%d", &u.health, &u.healthMax)
	fmt.Sscanf(line[2], "%d/%d", &u.magicka, &u.magickaMax)
	fmt.Sscanf(line[3], "%d/%d", &u.stamina, &u.staminaMax)
	fmt.Sscanf(line[4], "%d/%d", &u.ultimate, &u.ultimateMax)
	fmt.Sscanf(line[5], "%d/%d", &u.werewolf, &u.werewolfMax)
	// shield
	u.mapNX, _ = strconv.ParseFloat(line[7], 64)
	u.mapNY, _ = strconv.ParseFloat(line[8], 64)
	u.heading, _ = strconv.ParseFloat(line[9], 64)
	return u
}

func (u *UnitState) String() string {
	return fmt.Sprintf("%s (%d/%d)", u.UnitInfo.Name(), u.health, u.healthMax)
}

func (u *UnitState) ID() int {
	return u.unitID
}
