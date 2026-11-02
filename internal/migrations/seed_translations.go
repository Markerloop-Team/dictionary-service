package migrations

import (
	"embed"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:embed seeds/*.sql
var seedFiles embed.FS

// RunSeeds выполняет все SQL seed файлы при старте сервиса.
// Каждый файл использует ON CONFLICT DO UPDATE — безопасно запускать повторно.
func RunSeeds(db *gorm.DB, logger *logrus.Logger) error {
	entries, err := seedFiles.ReadDir("seeds")
	if err != nil {
		return fmt.Errorf("ошибка чтения директории seeds: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		data, err := seedFiles.ReadFile("seeds/" + entry.Name())
		if err != nil {
			logger.Warnf("Ошибка чтения seed файла %s: %v", entry.Name(), err)
			continue
		}

		sql := string(data)
		if strings.TrimSpace(sql) == "" {
			continue
		}

		if err := db.Exec(sql).Error; err != nil {
			logger.Warnf("Ошибка выполнения seed файла %s: %v", entry.Name(), err)
			continue
		}

		logger.Infof("Seed файл %s применён успешно", entry.Name())
	}

	return nil
}
