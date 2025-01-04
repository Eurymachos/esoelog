package esoelog

import (
	"testing"
)

func TestEnum(t *testing.T) {
	if lastLineType != EndTrial {
		t.Fatalf(`%d != %d`, lastLineType, EndTrial)
	}
}

/*
type lineTestCase struct {
	line     []string
	lineType LineType
	when     uint64
}

var testCases = []lineTestCase{
	{[]string{"4", "BEGIN_LOG", "1726530520043", "15", "NA Megaserver", "en", "eso.live.10.1"}, BeginLog, 4},
	{[]string{"4", "END_LOG"}, EndLog, 4},
	{[]string{"5565", "BEGIN_COMBAT"}, BeginCombat, 5565},
	{[]string{"7982", "END_COMBAT"}, EndCombat, 7982},
	{[]string{"5566", "PLAYER_INFO", "1",
		"[142210,26750,45557,45562,45048,45038,45053,45060,45071,45103,45084,45115,45150,45135,45155,29062,33293,39266,15594,30948,45509,13982,34741,203342,215493,63601]",
		"[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]",
		"[[HEAD,47971,F,46,ARMOR_TRAINING,MAGIC,43,MAGICKA,F,45,NORMAL],[NECK,34889,F,36,JEWELRY_ARCANE,ARCANE,98,MAGICKA_REGEN,F,36,ARCANE],[CHEST,46440,F,48,ARMOR_TRAINING,ARCANE,37,HEALTH,F,49,NORMAL],[SHOULDERS,7315,F,46,ARMOR_ORNATE,MAGIC,98,MAGICKA,F,46,MAGIC],[MAIN_HAND,46464,F,48,WEAPON_TRAINING,ARCANE,37,ABSORB_HEALTH,F,40,ARTIFACT],[WAIST,47975,F,46,ARMOR_TRAINING,MAGIC,43,HEALTH,F,43,NORMAL],[LEGS,47972,F,46,ARMOR_TRAINING,MAGIC,43,HEALTH,F,43,NORMAL],[FEET,47969,F,46,ARMOR_TRAINING,MAGIC,43,MAGICKA,F,42,NORMAL],[RING1,96769,T,4,JEWELRY_HEALTHY,MAGIC,107,HEALTH_REGEN,T,4,MAGIC],[RING2,27273,F,35,JEWELRY_ARCANE,ARCANE,0,REDUCE_SPELL_COST,F,35,ARCANE],[HAND,47970,F,46,ARMOR_TRAINING,MAGIC,43,HEALTH,F,45,NORMAL]]",
		"[25267,34835,25380,36028,36891,35460]",
		"[26768,25380,35414,36935]"}, PlayerInfo, 5566},
	/*{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	{[]string{}},
	// { []string{}, },
}

func TestLines(t *testing.T) {
	for i := range testCases {
		v, err := ParseLine(testCases[i].line)
		if err != nil {
			t.Fatalf(`ParseLine(%v) returned an unexpected error: %v`, testCases[i].line, err)
		}
		if v == nil {
			t.Fatalf(`ParseLine(%v) returned nil value but no error`, testCases[i].line)
		}
		if v.Type() != testCases[i].lineType {
			t.Fatalf(`ParseLine(%v) returned %v instead of %v`,
				testCases[i].line, v.Type(), testCases[i].lineType)
		}
	}
}
*/
