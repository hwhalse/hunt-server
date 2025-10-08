package handle_errors

import "log/slog"

func CheckError(err error, message string) {
	if err != nil {
		slog.Error(message, "error", err)
	}
}
