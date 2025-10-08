package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

func getZeroLogLoggingLevel() zerolog.Level {
	var lvl zerolog.Level
	switch os.Getenv("LOGGING_LEVEL") {
	case "PANIC":
		lvl = zerolog.PanicLevel
	case "FATAL":
		lvl = zerolog.FatalLevel
	case "ERROR":
		lvl = zerolog.ErrorLevel
	case "WARN":
		lvl = zerolog.WarnLevel
	case "INFO":
		lvl = zerolog.InfoLevel
	case "DEBUG":
		lvl = zerolog.DebugLevel
	case "TRACE":
		lvl = zerolog.TraceLevel
	default:
		lvl = zerolog.InfoLevel
	}
	return lvl
}

func NewLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	lvl := getZeroLogLoggingLevel()
	zerolog.SetGlobalLevel(lvl)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s |", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	logger := zerolog.New(output).With().Timestamp().Caller().Logger()
	return logger
}
