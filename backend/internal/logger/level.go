package logger

import (
	"log/slog"
)

type Level int

const (
	LevelInfo  = Level(10)
	LevelError = Level(11)
	LevelFatal = Level(12)
)

func (l Level) Level() slog.Level {
	return slog.Level(l)
}

func levelToString(l Level) string {
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
	lvl := Level(sloglevel.Level())
	a.Value = slog.StringValue(levelToString(lvl))
	return a
}
