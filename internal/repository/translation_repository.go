package repository

import (
	"context"
	"dictionary-service/internal/models"

	"gorm.io/gorm"
)

// TranslationRepository интерфейс для работы с репозиторием переводов
type TranslationRepository interface {
	// GetTranslations получает переводы по кодам, языку и приложению
	// Если app == nil, возвращаются общие переводы (app IS NULL) + переводы для указанного приложения
	// Если app != nil, возвращаются только переводы для указанного приложения
	GetTranslations(ctx context.Context, codes []string, language string, app *string) ([]models.TranslationDB, error)

	// GetAllTranslations получает все переводы для указанного языка и приложения
	GetAllTranslations(ctx context.Context, language string, app *string) ([]models.TranslationDB, error)

	// CreateTranslation создает новый перевод
	CreateTranslation(ctx context.Context, translation *models.TranslationDB) error

	// UpdateTranslation обновляет перевод
	UpdateTranslation(ctx context.Context, translation *models.TranslationDB) error

	// DeleteTranslation удаляет перевод
	DeleteTranslation(ctx context.Context, code, language string, app *string) error

	// GetTranslationsByCodes получает переводы для указанных кодов на всех языках
	// Возвращает map[code]map[language]value для удобного использования
	GetTranslationsByCodes(ctx context.Context, codes []string, app *string) (map[string]map[string]string, error)
}

// translationRepository реализация репозитория переводов
type translationRepository struct {
	db *gorm.DB
}

// CreateTranslationRepository создает репозиторий переводов
func CreateTranslationRepository(db *gorm.DB) TranslationRepository {
	return &translationRepository{db: db}
}

// GetTranslations получает переводы по кодам, языку и приложению
func (r *translationRepository) GetTranslations(ctx context.Context, codes []string, language string, app *string) ([]models.TranslationDB, error) {
	var translations []models.TranslationDB
	query := r.db.WithContext(ctx).Where("language = ?", language)

	// Фильтр по кодам
	if len(codes) > 0 {
		query = query.Where("code IN ?", codes)
	}

	// Фильтр по приложению
	if app == nil {
		// Если app не указан, возвращаем общие переводы (app IS NULL)
		query = query.Where("app IS NULL")
	} else {
		// Если app указан, возвращаем переводы для этого приложения ИЛИ общие (app IS NULL)
		// Приоритет: сначала специфичные для приложения, потом общие
		query = query.Where("app = ? OR app IS NULL", *app)
	}

	if err := query.Order("app DESC NULLS LAST").Find(&translations).Error; err != nil {
		return nil, err
	}

	// Убираем дубликаты, оставляя приоритет специфичным переводам для приложения
	if app != nil {
		translations = r.deduplicateTranslations(translations, *app)
	}

	return translations, nil
}

// GetAllTranslations получает все переводы для указанного языка и приложения
func (r *translationRepository) GetAllTranslations(ctx context.Context, language string, app *string) ([]models.TranslationDB, error) {
	var translations []models.TranslationDB
	query := r.db.WithContext(ctx).Where("language = ?", language)

	// Фильтр по приложению
	if app == nil {
		query = query.Where("app IS NULL")
	} else {
		query = query.Where("app = ? OR app IS NULL", *app)
	}

	if err := query.Order("app DESC NULLS LAST").Find(&translations).Error; err != nil {
		return nil, err
	}

	// Убираем дубликаты, оставляя приоритет специфичным переводам для приложения
	if app != nil {
		translations = r.deduplicateTranslations(translations, *app)
	}

	return translations, nil
}

// deduplicateTranslations убирает дубликаты переводов, оставляя приоритет специфичным для приложения
func (r *translationRepository) deduplicateTranslations(translations []models.TranslationDB, app string) []models.TranslationDB {
	seen := make(map[string]bool)
	result := make([]models.TranslationDB, 0, len(translations))

	// Сначала добавляем специфичные для приложения переводы
	for _, t := range translations {
		if t.App != nil && *t.App == app {
			key := t.Code + ":" + t.Language
			if !seen[key] {
				seen[key] = true
				result = append(result, t)
			}
		}
	}

	// Затем добавляем общие переводы (если специфичный не найден)
	for _, t := range translations {
		if t.App == nil {
			key := t.Code + ":" + t.Language
			if !seen[key] {
				seen[key] = true
				result = append(result, t)
			}
		}
	}

	return result
}

// CreateTranslation создает новый перевод
func (r *translationRepository) CreateTranslation(ctx context.Context, translation *models.TranslationDB) error {
	return r.db.WithContext(ctx).Create(translation).Error
}

// UpdateTranslation обновляет перевод
func (r *translationRepository) UpdateTranslation(ctx context.Context, translation *models.TranslationDB) error {
	updates := make(map[string]interface{})
	
	// Обновляем только изменяемые поля (исключаем ID, CreatedAt, DeletedAt)
	// Примечание: Code, Language, App обычно не меняются (они в unique index)
	// Но если нужно обновить, то обновляем их
	if translation.Code != "" {
		updates["code"] = translation.Code
	}
	if translation.Language != "" {
		updates["language"] = translation.Language
	}
	if translation.Value != "" {
		updates["value"] = translation.Value
	}
	if translation.App != nil {
		updates["app"] = translation.App
	} else {
		updates["app"] = nil
	}
	
	return r.db.WithContext(ctx).Model(translation).Updates(updates).Error
}

// DeleteTranslation удаляет перевод
func (r *translationRepository) DeleteTranslation(ctx context.Context, code, language string, app *string) error {
	query := r.db.WithContext(ctx).Where("code = ? AND language = ?", code, language)
	if app != nil {
		query = query.Where("app = ?", *app)
	} else {
		query = query.Where("app IS NULL")
	}
	return query.Delete(&models.TranslationDB{}).Error
}

// GetTranslationsByCodes получает переводы для указанных кодов на всех языках
// Возвращает map[code]map[language]value для удобного использования
func (r *translationRepository) GetTranslationsByCodes(ctx context.Context, codes []string, app *string) (map[string]map[string]string, error) {
	var translations []models.TranslationDB
	query := r.db.WithContext(ctx)

	// Фильтр по кодам
	if len(codes) > 0 {
		query = query.Where("code IN ?", codes)
	}

	// Фильтр по приложению
	if app == nil {
		query = query.Where("app IS NULL")
	} else {
		query = query.Where("app = ? OR app IS NULL", *app)
	}

	if err := query.Order("app DESC NULLS LAST").Find(&translations).Error; err != nil {
		return nil, err
	}

	// Убираем дубликаты, оставляя приоритет специфичным переводам для приложения
	if app != nil {
		translations = r.deduplicateTranslations(translations, *app)
	}

	// Группируем по коду и языку
	result := make(map[string]map[string]string)
	for _, t := range translations {
		if result[t.Code] == nil {
			result[t.Code] = make(map[string]string)
		}
		// Если уже есть значение для этого кода и языка, не перезаписываем (приоритет у специфичных)
		if _, exists := result[t.Code][t.Language]; !exists {
			result[t.Code][t.Language] = t.Value
		}
	}

	return result, nil
}
