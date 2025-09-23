package mng

import "github.com/makeevaas/project/sql-migrator/cfg"

type Migration struct {
	Up   string `yaml:"UP"`
	Down string `yaml:"DOWN"`
}

type Management struct {
	Cfg cfg.Config
}
