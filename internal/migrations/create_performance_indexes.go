package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// CreatePerformanceIndexes создает функциональные индексы для оптимизации поиска с LOWER()
// и обычные индексы на часто используемые поля
func CreatePerformanceIndexes(db *gorm.DB) error {
	// Проверяем существование индексов перед созданием
	indexes := []struct {
		name    string
		query   string
		checkQuery string
	}{
		// Функциональные индексы для поиска с LOWER()
		{
			name: "idx_countries_code_lower",
			query: `CREATE INDEX IF NOT EXISTS idx_countries_code_lower ON countries(LOWER(code))`,
			checkQuery: `SELECT 1 FROM pg_indexes WHERE indexname = 'idx_countries_code_lower'`,
		},
		{
			name: "idx_cities_code_lower",
			query: `CREATE INDEX IF NOT EXISTS idx_cities_code_lower ON cities(LOWER(code))`,
			checkQuery: `SELECT 1 FROM pg_indexes WHERE indexname = 'idx_cities_code_lower'`,
		},
		// Обычные индексы на часто используемые поля
		{
			name: "idx_translations_language_app",
			query: `CREATE INDEX IF NOT EXISTS idx_translations_language_app ON translations(language, app)`,
			checkQuery: `SELECT 1 FROM pg_indexes WHERE indexname = 'idx_translations_language_app'`,
		},
		{
			name: "idx_cities_country_code_created_at",
			query: `CREATE INDEX IF NOT EXISTS idx_cities_country_code_created_at ON cities(country_code, created_at)`,
			checkQuery: `SELECT 1 FROM pg_indexes WHERE indexname = 'idx_cities_country_code_created_at'`,
		},
	}

	for _, idx := range indexes {
		// Проверяем существование индекса
		var exists bool
		if err := db.Raw(idx.checkQuery).Scan(&exists).Error; err != nil {
			// Если таблица не существует, пропускаем
			continue
		}

		// Создаем индекс если его нет
		if err := db.Exec(idx.query).Error; err != nil {
			// Игнорируем ошибки если индекс уже существует или таблица не существует
			if !isIndexExistsError(err) {
				return fmt.Errorf("ошибка создания индекса %s: %w", idx.name, err)
			}
		}
	}

	return nil
}

// isIndexExistsError проверяет, является ли ошибка связанной с уже существующим индексом
func isIndexExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "already exists") || contains(errStr, "duplicate") || contains(errStr, "relation")
}

// contains проверяет наличие подстроки в строке (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsCaseInsensitive(s, substr))))
}

func containsCaseInsensitive(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}
