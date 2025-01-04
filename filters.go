package esoelog

// ZoneFilter filters a stream of log lines to only include lines that
// occur in a particular zone.
func ZoneFilter(zone string, in <-chan *LogLine, out chan<- *LogLine) {
	var inZone bool

	for l := range in {
		if l.LineType == ZoneChanged {
			inZone = l.LineData[3] == zone
		}
		if inZone || l.LineType == AbilityInfo {
			out <- l
		}
	}
	close(out)
}
