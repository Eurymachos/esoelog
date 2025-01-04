package esoelog

import (
	"fmt"
	"math"
)

func FindBoundary(msgs <-chan GameEvent) {
	var maxX, minX, maxY, minY float64
	var mapName string

	for msg := range msgs {
		switch m := msg.(type) {
		case *EventZoneEntered:
			mapName = m.Zone()
			maxX = -math.MaxFloat64
			minX = math.MaxFloat64
			maxY = -math.MaxFloat64
			minY = math.MaxFloat64

		case *EventZoneExited:
			if maxX != -math.MaxFloat64 &&
				minX != math.MaxFloat64 &&
				maxY != -math.MaxFloat64 &&
				minY != math.MaxFloat64 {
				fmt.Printf("%s X: [%f,%f] Y: [%f,%f]\n",
					mapName, minX, maxX, minY, maxY)
			} else {
				fmt.Printf("%s: insufficient information\n", mapName)
			}

		default:
			if us := m.Source(); us != nil {
				maxX = math.Max(maxX, us.mapNX)
				minX = math.Min(minX, us.mapNX)
				maxY = math.Max(maxY, us.mapNY)
				minY = math.Min(minY, us.mapNY)
			}
			if us := m.Target(); us != nil {
				maxX = math.Max(maxX, us.mapNX)
				minX = math.Min(minX, us.mapNX)
				maxY = math.Max(maxY, us.mapNY)
				minY = math.Min(minY, us.mapNY)
			}

		}
	}
}
