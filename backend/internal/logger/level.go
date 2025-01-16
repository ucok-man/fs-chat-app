package logger

import (
	"log/slog"
)

type level int

const (
	LevelInfo  = level(10)
	LevelError = level(11)
	LevelFatal = level(12)
)

func (l level) Level() slog.Level {
	return slog.Level(l)
}

func levelToString(l level) string {
	switch l {
	case 10:
		return "INFO"
	case 11:
		return "ERROR"
	case 12:
		return "FATAL"
	default:
		panic("invalid level value")
	}
}

func parselevel(a slog.Attr) slog.Attr {
	if !(a.Key == slog.LevelKey) {
		return a
	}

	sloglevel := a.Value.Any().(slog.Level)
	lvl := level(sloglevel.Level())
	a.Value = slog.StringValue(levelToString(lvl))
	return a
}
