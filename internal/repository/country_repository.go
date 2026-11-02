package repository

import (
	"context"
	"fmt"
	"strings"

	"dictionary-service/internal/models"

	"gorm.io/gorm"
)

// CountryRepository интерфейс для работы с репозиторием стран
type CountryRepository interface {
	// GetAllCountriesWithFilters получает страны с фильтрацией, сортировкой и пагинацией
	GetAllCountriesWithFilters(ctx context.Context, params *models.GetCountriesParams) ([]models.CountryDB, int64, error)

	// GetCountryByCode получает страну по коду
	GetCountryByCode(ctx context.Context, code string) (*models.CountryDB, error)

	// CreateCountry создает новую страну
	CreateCountry(ctx context.Context, country *models.CountryDB) error

	// UpdateCountry обновляет страну
	UpdateCountry(ctx context.Context, country *models.CountryDB) error

	// DeleteCountry удаляет страну
	DeleteCountry(ctx context.Context, code string) error
}

// countryRepository реализация репозитория стран
type countryRepository struct {
	db *gorm.DB
}

// CreateCountryRepository создает репозиторий стран
func CreateCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{db: db}
}

// GetAllCountriesWithFilters получает страны с фильтрацией, сортировкой и пагинацией
func (r *countryRepository) GetAllCountriesWithFilters(ctx context.Context, params *models.GetCountriesParams) ([]models.CountryDB, int64, error) {
	var countries []models.CountryDB
	var total int64

	query := r.db.WithContext(ctx).Model(&models.CountryDB{})

	// Применяем фильтры
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(code) LIKE ?", searchPattern)
	}

	// Применяем сортировку
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "code"
	}

	allowedSortFields := map[string]bool{
		"code": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "code"
	}

	sortOrder := strings.ToUpper(params.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	// Применяем пагинацию: cursor-based или offset-based
	// Cursor pagination работает только когда сортировка по ID
	useCursor := params.Cursor != "" && (sortBy == "id" || sortBy == "")
	if useCursor {
		// Cursor-based pagination: используем WHERE id > cursor вместо OFFSET
		// Конвертируем cursor (string) в uint для сравнения
		var cursorID uint
		if _, err := fmt.Sscanf(params.Cursor, "%d", &cursorID); err == nil {
			if sortOrder == "ASC" {
				query = query.Where("id > ?", cursorID)
			} else {
				query = query.Where("id < ?", cursorID)
			}
			// Для cursor pagination всегда сортируем по ID
			sortBy = "id"
		} else {
			// Если cursor невалидный, используем offset pagination
			useCursor = false
		}
	}
	if !useCursor {
		// Offset-based pagination (традиционная)
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset)
	}

	// Оптимизация: используем window function для получения total в одном запросе
	type CountryWithTotal struct {
		models.CountryDB
		Total int64 `gorm:"column:total"`
	}

	var results []CountryWithTotal
	if err := query.
		Select("id, code, created_at, updated_at, COUNT(*) OVER() as total").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(params.Limit).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	// Извлекаем countries и total из результатов
	countries = make([]models.CountryDB, len(results))
	if len(results) > 0 {
		total = results[0].Total
		for i, r := range results {
			countries[i] = r.CountryDB
		}
	} else {
		total = 0
	}

	return countries, total, nil
}

// GetCountryByCode получает страну по коду
func (r *countryRepository) GetCountryByCode(ctx context.Context, code string) (*models.CountryDB, error) {
	var country models.CountryDB
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&country).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

// CreateCountry создает новую страну
func (r *countryRepository) CreateCountry(ctx context.Context, country *models.CountryDB) error {
	return r.db.WithContext(ctx).Create(country).Error
}

// UpdateCountry обновляет страну
func (r *countryRepository) UpdateCountry(ctx context.Context, country *models.CountryDB) error {
	updates := make(map[string]interface{})
	
	// Обновляем только изменяемые поля (исключаем ID, CreatedAt, DeletedAt)
	// Примечание: Code обычно не меняется (он в unique index)
	// Но если нужно обновить, то обновляем его
	if country.Code != "" {
		updates["code"] = country.Code
	}
	
	return r.db.WithContext(ctx).Model(country).Updates(updates).Error
}

// DeleteCountry удаляет страну
func (r *countryRepository) DeleteCountry(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Where("code = ?", code).Delete(&models.CountryDB{}).Error
}
