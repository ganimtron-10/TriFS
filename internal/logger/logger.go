package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var globalLogger zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	globalLogger = zerolog.New(output).With().Timestamp().Logger()
}

func Log(component string, level zerolog.Level, msg string, args ...any) {

	formattedMsg := fmt.Sprintf("[%s] %s", component, msg)

	event := globalLogger.WithLevel(level)

	if len(args)%2 != 0 {
		event = event.Any("error", "odd number of arguments passed to logging.Log")
		event.Msg(formattedMsg)
		return
	}

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			event = event.Any("log_error", fmt.Sprintf("non-string key at index %d in logging.Log args", i))
			event.Msg(formattedMsg)
			return
		}
		value := args[i+1]
		event = event.Any(key, fmt.Sprintf("%#v", value))
	}

	event.Msg(formattedMsg)
}

func Debug(component string, msg string, args ...any) {
	Log(component, zerolog.DebugLevel, msg, args...)
}

func Info(component string, msg string, args ...any) {
	Log(component, zerolog.InfoLevel, msg, args...)
}

func Warn(component string, msg string, args ...any) {
	Log(component, zerolog.WarnLevel, msg, args...)
}

func Error(component string, msg string, args ...any) {
	Log(component, zerolog.ErrorLevel, msg, args...)
}
