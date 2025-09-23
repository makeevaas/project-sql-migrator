package cfg

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Ctx            context.Context
	Db             *pgx.Conn
	MigratePath    string
	MigrationFiles []string
}
