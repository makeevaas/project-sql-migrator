package mng

import "github.com/makeevaas/project-sql-migrator/pkg/cfg"

type Migration struct {
	Up   string `yaml:"up"`
	Down string `yaml:"down"`
}

type Management struct {
	Cfg cfg.Config
}

const GetMigrateDataReq = `SELECT version_id,is_applied,tstamp 
from db_version 
where version_id=$1 
ORDER BY tstamp 
DESC LIMIT 1;`
