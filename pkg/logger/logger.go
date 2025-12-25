package logger

import (
	"log/slog"

	"os"

	"github.com/lmittmann/tint"
)

func Init(level slog.Level) {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:     level,
		AddSource: true,
	})

	slog.SetDefault(slog.New(handler))
}
