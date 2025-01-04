
Unit IDs seem to be zone related rather than map (which isn't too much
of a surprise). I'm going to guess that unit related things such as
effect info and ability info persists with the units.

All units are removed when you leave a zone.
All units are removed when you leave the game.

The game seems to do a good job of keeping begin log and end log sane.

# Programs

## basics

Basics calculates certain basic statistics for a log file so that I
can verify certain assumptions about when various IDs are valid. e.g.

    Log file: Encounter-DR.log
    Abilities used undefined in zone 379885
    Abilities used undefined in log 0
    Map changes with 0 extant units 0
    Map changes with >0 extant units 79
    Zone changes with 0 extant units 34
    Zone changes with >0 extant units 0
    Log ends with 0 extant units 8
    Log ends with >0 extant units 0
    Log trailing lines 0

We can see that there were 79 map changes where units had been added,
but not removed while there were no zone changes like that. This leads
me to conclude that unit IDs remain valid across map changes, but not
zone changes.

We can also see that ability IDs seem to be valid for the entire log
rather than for just the zone because we can see 379885 ability uses
in zones where there was no ability info for that ability in that
zone. On the other hand, every ability use is preceded by ability info
in that log.

## faolchu

Faolchu prints out the encounter log for an entire zone in a format
that is slightly higher level and easier to read than the raw
encounter log. The default zone it selects is Camlorn Keep, hence the
name. Other zones can be selected with a `-z ZONENAME` command line
option.

I use this program to investigate undocumented fight mechanics.

## limits

Limits prints out the X and Y limits of all the zones in the encounter
log. It's not very useful.

## monsters

Monsters prints out a list of all hostile monsters you have seen along
with the abilities you have seen them use. It also prints out the
maximum health of each monster when that information is available.

## split

Split splits an encounter log into it's component log files. ESO
creates a new log in the file every time you log in with a
character. Split puts each log into it's own file with a date and time
in the file name.

I use this when I've been collecting log data for a week and then I
want to see what happened in this morning's session without processing
all the other data.

The log files that are created are written out as standard CSV files
and not in ESO's ad-hoc, broken CSV format. This means that you may
not be able to upload these output files into esologs.com. Keep your
original encounter logs.

# BUGS

Sometimes the encounter log lists a player's pet as HOSTILE instead of
NPC_ALLY and so we report it as a monster, which ends up looking a bi
strange.
