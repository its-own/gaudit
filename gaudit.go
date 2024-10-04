package gaudit

import (
	_ "github.com/its-own/gaudit/internal/audit_log"
	"log/slog"
)

func init() {
	slog.Default().Info("hookiee in action")
}
