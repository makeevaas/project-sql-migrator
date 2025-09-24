package mng

import "github.com/makeevaas/project/sql-migrator/pkg/cfg"

type Migration struct {
	Up   string `yaml:"up"`
	Down string `yaml:"down"`
}

type Management struct {
	Cfg cfg.Config
}
