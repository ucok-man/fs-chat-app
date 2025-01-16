package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type option func(l *Logger)

func WithLevel(lvl level) option {
	return func(l *Logger) {
		l.minlvl = lvl
	}
}

type Logger struct {
	logger *slog.Logger
	minlvl level
	ctx    struct {
		level level
		attrs []slog.Attr
		msg   string
	}
}

func New(opts ...option) *Logger {
	l := &Logger{
		minlvl: LevelInfo,
	}

	for _, optFn := range opts {
		optFn(l)
	}

	handleopt := &slog.HandlerOptions{
		Level:     l.minlvl,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			a = parselevel(a)
			a = parsetrace(a)
			return a
		},
	}

	l.logger = slog.New(slog.NewJSONHandler(os.Stdout, handleopt))
	return l

}

func (l *Logger) resetCtx() {
	l.ctx.attrs = []slog.Attr{}
	l.ctx.msg = ""
}

func (l *Logger) Info(format string, a ...any) *Logger {
	l.ctx.msg = fmt.Sprintf(format, a...)
	l.ctx.level = LevelInfo
	return l
}

func (l *Logger) Error(err error) *Logger {
	l.ctx.msg = err.Error()
	l.ctx.level = LevelError
	l.Attr("trace", err)
	return l
}

func (l *Logger) Fatal(err error) *Logger {
	l.ctx.msg = err.Error()
	l.ctx.level = LevelFatal
	l.Attr("trace", err)
	return l
}

func (l *Logger) Attr(key string, value any) *Logger {
	l.ctx.attrs = append(l.ctx.attrs, slog.Any(key, value))
	return l
}

func (l *Logger) Send() {
	l.logger.LogAttrs(context.Background(), l.ctx.level.Level(), l.ctx.msg, l.ctx.attrs...)
	l.resetCtx()
}
