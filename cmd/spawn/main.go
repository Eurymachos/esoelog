package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	//"github.com/lukegb/dds"

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
		//events := make(chan eso.GameEvent)

		go eso.LogReader(fn, log)
		go eso.ZoneFilter(zone, log, filtered)
		//go eso.RunGame(filtered, events)
		//PrintFight(events)
		PrintSpawns(filtered)
	}

}

// We want to print where on the map a unit spawns.
// We need to know which map we are one.
// We need to know unit location, but only first location for each map.
func PrintSpawns(c <-chan *eso.LogLine) {
	allUnits := make(map[int]*eso.UnitInfo)
	unitFound := make(map[int]bool)

	var gameMap *image.RGBA

	for ll := range c {
		switch ll.LineType {
		case eso.MapChanged:
			fmt.Println(ll.LineData[3], ll.LineData[4])
			gameMap = loadMap(ll.LineData[4])

		case eso.UnitAdded:
			u := eso.NewUnitInfo(ll.LineData[2:])
			allUnits[u.ID()] = u

		case eso.CombatEvent:
			src, tgt := getUnits(ll.LineData[9:])
			if src != nil && !unitFound[src.ID()] {
				fmt.Println("Unit", src.ID(),
					src.X(), src.Y(), src.D())
				ping(gameMap, src.X(), src.Y())
				unitFound[src.ID()] = true
			}
			if tgt != nil && !unitFound[tgt.ID()] {
				fmt.Println("Unit", tgt.ID(),
					tgt.X(), tgt.Y(), tgt.D())
				unitFound[tgt.ID()] = true
			}

			/*
				__BEGIN_CAST__ - durationMS, channeled, castTrackId, abilityId, _sourceUnitState_, _targetUnitState_
				__COMBAT_EVENT__ - actionResult, damageType, powerType, hitValue, overflow, castTrackId, abilityId, _sourceUnitState_, _targetUnitState_
								__HEALTH_REGEN__ - effectiveRegen, _unitState_

				__UNIT_ADDED__ - unitId, unitType, isLocalPlayer, playerPerSessionId, monsterId, isBoss, classId, raceId, name, displayName, characterId, level, championPoints, ownerUnitId, reaction, isGroupedWithLocalPlayer
				__UNIT_CHANGED__ - unitId, classId, raceId, name, displayName, characterId, level, championPoints, ownerUnitId, reaction, isGroupedWithLocalPlayer
				__EFFECT_CHANGED__ - changeType, stackCount, castTrackId, abilityId, _sourceUnitState_, _targetUnitState_, playerInitiatedRemoveCastTrackId:optional

								__MAP_INFO__ - id, name, texturePath
			*/
		}
	}

	w, err := os.Create("out.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	png.Encode(w, gameMap)
	w.Close()
}

func ping(img *image.RGBA, x, y float64) {
	r := img.Bounds()
	xi := int(float64(r.Max.X) * x)
	yi := int(float64(r.Max.Y) * y)
	red := color.RGBA{255, 0, 0, 255}
	img.Set(xi+1, yi+1, red)
	img.Set(xi+1, yi, red)
	img.Set(xi+1, yi-1, red)
	img.Set(xi, yi+1, red)
	img.Set(xi, yi, red)
	img.Set(xi, yi-1, red)
	img.Set(xi-1, yi+1, red)
	img.Set(xi-1, yi, red)
	img.Set(xi-1, yi-1, red)
}

func getUnits(line []string) (src, dst *eso.UnitState) {
	src = eso.NewUnitState(line[:10])
	dst = src
	if line[10] != "*" {
		dst = eso.NewUnitState(line[10:20])
	}
	return
}

func loadMap(basename string) *image.RGBA {
	var aMap *image.RGBA

	for dp, i := image.ZP, 0; i < 9; {
		filename := fmt.Sprintf("Documents/ESO/game/art/maps/%s_%d.png", basename, i)
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		img, err := png.Decode(f)
		if err != nil {
			f.Close()
			fmt.Println(err)
			return nil
		}
		f.Close()

		sr := img.Bounds()
		if i == 0 {
			r := image.Rectangle{image.ZP, sr.Max.Mul(3)}
			aMap = image.NewRGBA(r)
		}

		r := image.Rectangle{dp, dp.Add(sr.Size())}
		draw.Draw(aMap, r, img, sr.Min, draw.Src)

		i += 1
		dp.X += 256
		if i%3 == 0 {
			dp.Y += 256
			dp.X = 0
		}
	}

	return aMap
}
