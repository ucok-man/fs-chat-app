package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type option func(l *Logger)

func WithLevel(lvl Level) option {
	return func(l *Logger) {
		l.minlvl = lvl
	}
}

type Logger struct {
	logger *slog.Logger
	minlvl Level
	ctx    struct {
		skipC int
		trace slog.Attr
		level Level
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
			// parse source
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}

	l.logger = slog.New(slog.NewJSONHandler(os.Stdout, handleopt))
	return l

}

func (l *Logger) SetDefault() {
	slog.SetDefault(l.logger)
}

func (l *Logger) resetCtx() {
	l.ctx.attrs = []slog.Attr{}
	l.ctx.msg = ""
	l.ctx.trace = slog.Attr{}
	l.ctx.skipC = 0
}

func (l *Logger) Info(format string, a ...any) *Logger {
	l.ctx.msg = fmt.Sprintf(format, a...)
	l.ctx.level = LevelInfo
	return l
}

func (l *Logger) Error(err error) *Logger {
	l.ctx.msg = err.Error()
	l.ctx.level = LevelError
	l.ctx.trace = slog.Any("trace", err)
	return l
}

func (l *Logger) Fatal(err error) *Logger {
	l.ctx.msg = err.Error()
	l.ctx.level = LevelFatal
	l.ctx.trace = slog.Any("trace", err)
	return l
}

func (l *Logger) Attr(key string, value any) *Logger {
	l.ctx.attrs = append(l.ctx.attrs, slog.Any(key, value))
	return l
}

func (l *Logger) SkipC(skip int) *Logger {
	l.ctx.skipC = skip
	return l
}

func (l *Logger) Send() {
	var pcs [1]uintptr
	runtime.Callers(2+l.ctx.skipC, pcs[:])
	l.ctx.attrs = append(l.ctx.attrs, l.ctx.trace)
	record := slog.NewRecord(time.Now(), l.ctx.level.Level(), l.ctx.msg, pcs[0])
	l.logger.Handler().WithAttrs(l.ctx.attrs).Handle(context.Background(), record)
	l.resetCtx()
}
