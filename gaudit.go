package gaudit

import (
	"github.com/its-own/gaudit/in"
	_ "github.com/its-own/gaudit/internal/audit_log"
	"github.com/its-own/gaudit/internal/hooks"
	amgo "github.com/its-own/gaudit/internal/infracture/db/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

func init() {
	slog.Default().Info("gaudit in action")
}

type Config struct {
	*mongo.Client
	Database *mongo.Database
	hook     in.Hook
	Logger   *slog.Logger
}

func Init(c *Config) *amgo.Mongo {
	return amgo.InitMongo(c.Client, c.Database, hooks.NewDefaultHook(c.Logger))
}
