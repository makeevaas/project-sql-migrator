package cfg

import (
	"context"

	"github.com/makeevaas/project/sql-migrator/pkg/db"
)

type Config struct {
	Ctx            context.Context
	DB             *db.DB
	MigratePath    string
	MigrationFiles []string
}
