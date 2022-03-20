package tok

import (
	"fmt"
	"strconv"
	"strings"
)

type LogEntry struct {
	EnterAt int
	Info    string
	Level   int
	ExitAt  int
	Error   error
}

func (e LogEntry) String() string {
	return fmt.Sprintf("%d%s@ %d %s", e.Level, strings.Repeat(".", e.Level), e.EnterAt, e.Info)
}

type Log struct {
	Entries []LogEntry
	level   int
}

func MonitorGrammar(g Grammar) *Log {
	l := &Log{}
	rules := CollectRuleReaders(g)
	l.Monitor(rules...)
	return l
}

func (l *Log) Enter(info string, pos int) *LogEntry {
	l.level++
	i := len(l.Entries)
	l.Entries = append(l.Entries, LogEntry{
		EnterAt: pos,
		Info:    info,
		Level:   l.level,
	})
	return &l.Entries[i]
}

func (l *Log) Exit(e *LogEntry, pos int, err error) {
	e.ExitAt = pos
	e.Error = err
	l.level--
}

func (l *Log) Monitor(readers ...*RuleReader) {
	for _, r := range readers {
		r.Monitor(l)
	}
}

func (l *Log) Print() {
	for _, e := range l.Entries {
		fmt.Println(e)
	}
}

func (l *Log) PrintWithPreview(str string, n int) {
	for _, e := range l.Entries {
		fmt.Println(e, ">", strconv.Quote(subStringFrom(str, e.EnterAt, n)))
	}
}

func (l *Log) Reset() {
	l.Entries = []LogEntry{}
	l.level = 0
}

type monitorReader struct {
	info string
	log  *Log
	sub  Reader
}

func (r *monitorReader) Read(s *Scanner) error {
	entry := r.log.Enter(r.info, s.pos)
	err := r.sub.Read(s)
	r.log.Exit(entry, s.pos, err)
	return err
}

func (r *monitorReader) What() string {
	return r.sub.What()
}

func Monitor(r Reader, l *Log, info string) Reader {
	return &monitorReader{
		info: info,
		log:  l,
		sub:  r,
	}
}
