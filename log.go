package esoelog

//go:generate stringer -type=LineType

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/iancoleman/strcase"
)

type LineType int

const (
	BeginLog LineType = iota
	EndLog
	BeginCombat
	EndCombat
	PlayerInfo
	BeginCast
	EndCast
	CombatEvent
	HealthRegen
	UnitAdded
	UnitChanged
	UnitRemoved
	EffectChanged
	AbilityInfo
	EffectInfo
	MapChanged  // MapInfo in the docs
	ZoneChanged // ZoneInfo in the docs
	TrialInit
	BeginTrial
	EndTrial
	lastLineType           = iota - 1
	firstLineType LineType = 0
)

type LogLine struct {
	LineType
	DeltaMs  time.Duration
	LineData []string
}

var lineTags map[string]LineType

func init() {
	lineTags = make(map[string]LineType, lastLineType+1)

	for i := firstLineType; i <= lastLineType; i++ {
		tag := strcase.ToScreamingSnake(i.String())
		lineTags[tag] = i
	}
}

// LogReader opens and reads an ESO encounter log.
//
// An ESO encounter log is a CSV file that is mostly compatible with
// [encoding/csv].
func LogReader(filename string, lines chan<- *LogLine) {
	defer close(lines)

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	line, err := r.Read()
	for ; err == nil; line, err = robustRead(r) {
		ms, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		if lt, found := lineTags[line[1]]; found {
			log := &LogLine{
				DeltaMs:  time.Millisecond * time.Duration(ms),
				LineType: lt,
				LineData: line,
			}
			lines <- log
		} else {
			fmt.Println("Tag", line[1], "not found:", line)
		}
	}
	if err != io.EOF {
		fmt.Println(err)
		err = nil
	}
}

func robustRead(r *csv.Reader) ([]string, error) {
	rval, err := r.Read()
	if perr, ok := err.(*csv.ParseError); ok {
		if perr.Unwrap() == csv.ErrQuote {
			// Sigh
			err = nil
			rval = append(rval, []string{`"Eyeballs"`, "/esoui/art/icons/event_halloween_2016_skull_cup_grapes.dds", "T", "T"}...)
		}
	}

	return rval, err
}

// LogSplitter splits the input data into separate log files for each
// game session. This makes it easier to select data for further
// processing.
//
// The output data is in Go's CSV format rather than ESO's CSV format
// and so the output from this may not be compatible with esologs.com.
func LogSplitter(c <-chan *LogLine) {
	var f *os.File
	var w *csv.Writer

	for line := range c {
		if line.LineType == BeginLog {
			epochMs, err := strconv.ParseInt(line.LineData[2], 10, 64)
			if err != nil {
				fmt.Println(err)
			}

			t := time.UnixMilli(epochMs)
			tt := t.Format("Encounter_2006-01-02@15_04_05.log")
			fmt.Println(line.DeltaMs, line.LineType, tt)
			f, err = os.Create(tt)
			if err != nil {
				fmt.Println(err)
			} else {
				defer f.Close()
			}

			w = csv.NewWriter(f)
			defer w.Flush()
		}

		err := w.Write(line.LineData)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
