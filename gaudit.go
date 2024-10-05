package gaudit

import (
	_ "github.com/its-own/gaudit/internal/audit_log"
	"github.com/its-own/gaudit/internal/hooks"
	"log/slog"
)

func init() {
	slog.Default().Info("gaudit in action")
}

func New() *hooks.DefaultHooks {
	return hooks.NewDefaultHook()
}
