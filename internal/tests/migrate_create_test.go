package tests

import (
	"context"
	"os"
	"testing"

	mng "github.com/makeevaas/project-sql-migrator/internal"
	"github.com/makeevaas/project-sql-migrator/pkg/cfg"
)

// TestAdd проверяет функцию Add.
func TestCreateFileMigration(t *testing.T) {
	m := &mng.Management{
		Cfg: cfg.Config{Ctx: context.Background(), MigratePath: "../../migrations"},
	}
	file, err := m.CreateFileMigration()
	if err != nil {
		t.Fatalf("Не удалось создать файл: %v", err)
	}
	t.Log("Создан временный файл - ", file)

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		t.Fatalf("Файл %s не найден", file)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Не удалось прочитать файл: %v", err)
	}

	comment := `# The SQL in the UP section is executed when applying a migration.
# The SQL in the DOWN section is executed when rolling back a migration.
up: ""
down: ""
`
	if string(data) != comment {
		t.Errorf("Содержимое файла не совпадает. Ожидалось: \n%s, получено: \n%s", comment, string(data))
	}

	os.Remove(file)
	t.Log("Временный файл удален - ", file)
}
