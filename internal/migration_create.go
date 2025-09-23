package mng

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

func (m *Management) CreateFileMigration() error {
	fileName := "/" + time.Now().Format("20060102150405") + "_.yaml"
	file, err := os.Create(m.Cfg.MigratePath + fileName)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	defer file.Close()

	// Записываем данные в файл
	config := &Migration{}
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("error yaml marshal: %w", err)
	}
	comment := `# The SQL in the UP section is executed when applying a migration.
# The SQL in the DOWN section is executed when rolling back a migration.
`
	_, err = file.WriteString(comment)
	if err != nil {
		return fmt.Errorf("error write comment string: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	return nil
}
