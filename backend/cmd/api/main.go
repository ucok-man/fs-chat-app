package main

import (
	"fmt"

	"github.com/ucok-man/fs-chat-app-backend/internal/logger"
)

// type application struct {
// }

func main() {
	// logger := logger.New(logger.WithJSON(), logger.WithLevel())
	// slog.
	l := logger.New(logger.WithLevel(logger.LevelError))
	l.Info("Hello World").Attr("some", "value").Send()
	l.Error(fmt.Errorf("WTF Error")).Attr("hello", "world").Send()
	l.Fatal(fmt.Errorf("FATAL ERROR")).Attr("x", "y").Send()
}
