// Package utils contains small cross-cutting helpers.
package utils

import "log/slog"

// HandleErrorOrLogWithMessages logs either an error message or success message.
func HandleErrorOrLogWithMessages(err error, errMsg string, successMsg string) {
	if err != nil {
		slog.Error(errMsg, "error", err)
		return
	}
	if successMsg != "" {
		slog.Info(successMsg)
	}
}
