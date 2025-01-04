package esoelog

// GameEvent implements a higher level concept than the original log lines.
type GameEvent interface {
	Source() *UnitState
	Target() *UnitState
	Zone() string
}

type eventZone string

func (e eventZone) Zone() string {
	return string(e)
}

// EventZoneEntered is generated whenever the player enters a new zone.
type EventZoneEntered struct {
	eventZone
}

func (e *EventZoneEntered) Source() *UnitState {
	return nil
}

func (e *EventZoneEntered) Target() *UnitState {
	return nil
}

// NewEventZoneEntered creates a new ZoneEntered event for the named zone.
func NewEventZoneEntered(zoneName string) *EventZoneEntered {
	return &EventZoneEntered{eventZone: eventZone(zoneName)}
}

// EventZoneExited is generated whenever the player leaves a map (zone).
type EventZoneExited struct {
	eventZone
}

// NewEventZoneExited creates a new ZoneExited event for the named zone.
func NewEventZoneExited(zoneName string) *EventZoneExited {
	return &EventZoneExited{eventZone: eventZone(zoneName)}
}

func (e *EventZoneExited) Source() *UnitState {
	return nil
}

func (e *EventZoneExited) Target() *UnitState {
	return nil
}

// EventMap is generated whenever there is a map transition.
// Unlike with zones, there is no event for leaving a map because they
// just aren't that important.
// The player cares about which map they are currently in, but the
// game doesn't.
type EventMap struct {
	eventZone
	MapName string // Often the same as the zone name
	MapArt  string // The filename of the ingame map (unique?)
}

// NewEventMap creates a new EventMap for a map transition.
func NewEventMap(zoneName, mapName, mapArt string) *EventMap {
	return &EventMap{
		eventZone: eventZone(zoneName),
		MapName:   mapName,
		MapArt:    mapArt,
	}
}

func (e *EventMap) Source() *UnitState {
	return nil
}

func (e *EventMap) Target() *UnitState {
	return nil
}

type EventCombat struct {
	ActionResult string
	DamageType   string
	powerType    string
	HitValue     string
	overflow     string
	Ability      string
	When         string
	src          *UnitState
	tgt          *UnitState
	eventZone

	// castTrackID
}

func (e *EventCombat) Source() *UnitState {
	// A quick validity check
	if e.src != nil && e.src.UnitInfo != nil {
		return e.src
	}
	return nil
}

func (e *EventCombat) Target() *UnitState {
	// A quick validity check
	if e.tgt != nil && e.tgt.UnitInfo != nil {
		return e.tgt
	}
	return nil
}

type EventEffect struct {
	ChangeType string
	Ability    string
	When       string
	src        *UnitState
	tgt        *UnitState
	eventZone
	// castTrackID
}

func (e *EventEffect) Source() *UnitState {
	// A quick validity check
	if e.src != nil && e.src.UnitInfo != nil {
		return e.src
	}
	return nil
}

func (e *EventEffect) Target() *UnitState {
	// A quick validity check
	if e.tgt != nil && e.tgt.UnitInfo != nil {
		return e.tgt
	}
	return nil
}

// EventUnitSeen is a synthetic event used to list all the units seen
// in a zone.
// One event is created for each unit seen.
type EventUnitSeen struct {
	*UnitState
	eventZone
}

func NewEventUnitSeen(zoneName string, us *UnitState) *EventUnitSeen {
	return &EventUnitSeen{us, eventZone(zoneName)}
}

func (e *EventUnitSeen) Source() *UnitState {
	return e.UnitState
}

func (e *EventUnitSeen) Target() *UnitState {
	return e.UnitState
}

// BeginCombat
// EndCombat
// UnitInfo - at end, consolidation of UnitAdded, UnitChanged and UnitRemoved
// also unit state info from other records

// CombatEvent
// Cast (summary of BeginCast & EndCast)
// Effect

// MapExited

// What about

// BeginLog ignored?
// EndLog ignored?
// PlayerInfo reported at map end?
// CombatEvent
// HealthRegen
// EffectChanged
// AbilityInfo
// EffectInfo

// MapChanged  // MapInfo in the docs
// ZoneChanged // ZoneInfo in the docs
// TrialInit
// BeginTrial
// EndTrial
