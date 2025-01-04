package esoelog

import (
	"fmt"
	"strconv"
	"time"
)

type gameState struct {
	logStart    time.Time
	combatStart time.Time
	zoneID      int
	zoneName    string
	units       map[int]*UnitInfo
	unitStates  map[int]*UnitState
	abilities   map[int]string
	casting     map[int]bool
}

func NewGameState() *gameState {
	g := new(gameState)
	g.reset()

	return g
}

func (g *gameState) reset() {
	g.units = make(map[int]*UnitInfo)
	g.unitStates = make(map[int]*UnitState)
	if g.abilities == nil {
		g.abilities = make(map[int]string)
	}
	g.casting = make(map[int]bool)
}

func RunGame(c <-chan *LogLine, m chan<- GameEvent) {
	g := NewGameState()

	for ll := range c {
		switch ll.LineType {
		case BeginLog:
			g.doBeginLog(ll.DeltaMs, ll.LineData)
		case BeginCombat:
			g.doBeginCombat(ll.DeltaMs, ll.LineData)
		case EndCombat:
			g.doEndCombat(ll.DeltaMs, ll.LineData)
		case UnitAdded:
			g.doUnitAdded(ll.DeltaMs, ll.LineData)
		case UnitRemoved:
			g.doUnitRemoved(ll.DeltaMs, ll.LineData)
		case MapChanged: // do nothing
		case ZoneChanged:
			if g.zoneName != "" {
				m <- NewEventZoneExited(g.zoneName)
			}
			for _, us := range g.unitStates {
				m <- NewEventUnitSeen(g.zoneName, us)
			}
			g.doZoneChanged(ll.DeltaMs, ll.LineData)
			m <- NewEventZoneEntered(g.zoneName)
		case CombatEvent:
			m <- g.doCombatEvent(ll.DeltaMs, ll.LineData)
		case AbilityInfo:
			g.doAbilityInfo(ll.DeltaMs, ll.LineData)
		case EffectChanged:
			m <- g.doEffectChanged(ll.DeltaMs, ll.LineData)
		case BeginCast:
			g.doBeginCast(ll.DeltaMs, ll.LineData)
		case EndCast:
			g.doEndCast(ll.DeltaMs, ll.LineData)
		case HealthRegen:
		case PlayerInfo:
		case EffectInfo:
		case UnitChanged:
			g.doUnitChanged(ll.DeltaMs, ll.LineData)
		case EndLog:
		default:
			fmt.Println(ll.LineType, ll.LineData)
		}
	}
	// FIXME
	m <- NewEventZoneExited(g.zoneName)
	for _, us := range g.unitStates {
		m <- NewEventUnitSeen(g.zoneName, us)
	}
	close(m)
}

func (g *gameState) Time(when time.Duration) string {
	now := g.logStart.Add(when)
	combatTime := now.Sub(g.combatStart)
	if combatTime.Minutes() < 30 {
		return combatTime.String()
	}

	return now.Format(time.TimeOnly)
}

// BEGIN_LOG, timeSinceEpocsMS, logVersion, realmName, language, gameVersion
func (g *gameState) doBeginLog(when time.Duration, line []string) {
	start, err := strconv.ParseInt(line[2], 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	g.logStart = time.UnixMilli(start)
	//	fmt.Println("Log started at", g.logStart)
	//	fmt.Println("First entry at", g.logStart.Add(when))
}

// UNIT_ADDED, unitId, unitType, isLocalPlayer, playerPerSessionId,
// monsterId, isBoss, classId, raceId, name, displayName, characterId,
// level, championPoints, ownerUnitId, reaction,
// isGroupedWithLocalPlayer
//
// 4520901,UNIT_ADDED,5,MONSTER,F,0,15134,F,0,0,"Faolchu the
// Reborn","",0,50,160,0,HOSTILE,F
func (g *gameState) doUnitAdded(when time.Duration, line []string) {

	u := NewUnitInfo(line[2:])

	g.units[u.unitID] = u

	/*if g.location == CamlornKeep {
		fmt.Println("Unit:", u.name, "Lvl", u.level, "CP", u.championPoints)
	}*/
}

// UNIT_CHANGED - unitId, classId, raceId, name, displayName,
// characterId, level, championPoints, ownerUnitId, reaction,
// isGroupedWithLocalPlayer
func (g *gameState) doUnitChanged(when time.Duration, line []string) {
	unitID, _ := strconv.Atoi(line[2])
	u := g.units[unitID]
	u.name = line[5]
	u.displayName = line[6]
	u.Reaction = line[11]
	//fmt.Println(u.name, "-", unitID, "became", u.Reaction)
}

// UNIT_REMOVED, unitId
func (g *gameState) doUnitRemoved(when time.Duration, line []string) {
	unitID, _ := strconv.Atoi(line[2])
	delete(g.units, unitID)
}

// ZONE_CHANGED, id, name, dungeonDifficulty
func (g *gameState) doZoneChanged(when time.Duration, line []string) {
	id, err := strconv.Atoi(line[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	g.reset()
	g.zoneID = id
	g.zoneName = line[3]
	if line[4] != "NONE" {
		g.zoneName += fmt.Sprintf(" (%s)", line[4])
	}

	return
}

// COMBAT_EVENT, actionResult, damageType, powerType, hitValue,
// overflow, castTrackId, abilityId, sourceUnitState, targetUnitState
func (g *gameState) doCombatEvent(when time.Duration, line []string) *EventCombat {
	e := new(EventCombat)
	e.ActionResult = line[2]
	e.DamageType = line[3]
	e.HitValue = line[5]
	//castTrackId := line[7]
	abilityID, _ := strconv.Atoi(line[8])
	e.Ability = g.abilities[abilityID]
	e.When = g.Time(when)

	e.eventZone = eventZone(g.zoneName)
	e.src, e.tgt = g.getUnits(line[9:])

	return e
}

func (g *gameState) doBeginCombat(when time.Duration, line []string) {
	g.combatStart = g.logStart.Add(when)
}

func (g *gameState) doEndCombat(when time.Duration, line []string) {
	g.combatStart = time.Time{}
}

// ABILITY_INFO, abilityId, name, iconPath, interruptible, blockable
// [, effect1, effect2, effect3]
// Optional three extra effects for scribed skills?
func (g *gameState) doAbilityInfo(when time.Duration, line []string) {

	abilityID, _ := strconv.Atoi(line[2])
	if len(line) == 3 && line[2] == "84700" {
		g.abilities[abilityID] = `"Eyeballs"` // Thanks ESO
		return
	}

	g.abilities[abilityID] = line[3]
}

// EFFECT_CHANGED, changeType, stackCount, castTrackId, abilityId,
// _sourceUnitState_, _targetUnitState_,
// playerInitiatedRemoveCastTrackId:optional
func (g *gameState) doEffectChanged(when time.Duration, line []string) *EventEffect {
	//fmt.Println(line)
	e := new(EventEffect)
	e.ChangeType = line[2]
	//e.castTrackId, _ := strconv.Atoi(line[4])
	abilityID, _ := strconv.Atoi(line[5])
	e.Ability = g.abilities[abilityID]

	e.When = g.Time(when)

	e.src, e.tgt = g.getUnits(line[6:])
	e.eventZone = eventZone(g.zoneName)

	if e.Ability == "" {
		fmt.Println(line)
	}

	return e
}

// BEGIN_CAST, durationMS, channeled, castTrackId, abilityId,
// _sourceUnitState_, _targetUnitState_
func (g *gameState) doBeginCast(when time.Duration, line []string) {
	castTrackID, _ := strconv.Atoi(line[2])
	g.casting[castTrackID] = true
	// if g.location == CamlornKeep {
	// 	fmt.Println("Casting", castTrackId, line)
	// }
}

// END_CAST, endReason, castTrackId, interruptingAbilityId:optional,
// interruptingUnitId:optional
func (g *gameState) doEndCast(when time.Duration, line []string) {
	//	castTrackId, _ := strconv.Atoi(line[1])
	//delete(g.casting, castTrackId)
	// if g.location == CamlornKeep {
	// 	fmt.Println("End cast", castTrackId, line[0])
	// }
}

func (g *gameState) getUnits(line []string) (src, dst *UnitState) {
	src = parseUnitState(line[:10])
	if src != nil {
		BindUnit(g.units[src.unitID], src)
		g.unitStates[src.unitID] = src
	}
	dst = src
	if line[10] != "*" {
		dst = parseUnitState(line[10:20])
		if dst != nil {
			BindUnit(g.units[dst.unitID], dst)
			g.unitStates[dst.unitID] = dst
		}
	}
	return
}

func BindUnit(ui *UnitInfo, us *UnitState) {
	us.UnitInfo = ui
	for _, i := range ui.maxHealth {
		if i == us.healthMax {
			return
		}
	}
	ui.maxHealth = append(ui.maxHealth, us.healthMax)
}
