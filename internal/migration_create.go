package mng

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func (m *Management) CreateFileMigration() error {
	fileName := "/" + time.Now().Format("20060102150405") + "_.yaml"
	file, err := os.Create(m.Cfg.MigratePath + fileName)
	if err != nil {
		log.Fatalf("failed to create migration file: %v", err)
	}
	defer file.Close()

	// Записываем данные в файл
	config := &Migration{}
	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error yaml marshal: %v", err)
	}
	comment := `# The SQL in the UP section is executed when applying a migration.
# The SQL in the DOWN section is executed when rolling back a migration.
`
	_, err = file.WriteString(comment)
	if err != nil {
		log.Fatalf("error write comment string: %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("failed to write migration file: %v", err)
	}

	return nil
}
