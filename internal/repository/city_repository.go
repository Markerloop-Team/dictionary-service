package repository

import (
	"context"
	"fmt"
	"strings"

	"dictionary-service/internal/models"

	"gorm.io/gorm"
)

// CityRepository интерфейс для работы с репозиторием городов
type CityRepository interface {
	// GetCitiesByCountryCode получает города по коду страны
	GetCitiesByCountryCode(ctx context.Context, countryCode string) ([]models.CityDB, error)

	// GetCitiesWithFilters получает города с фильтрацией, сортировкой и пагинацией
	GetCitiesWithFilters(ctx context.Context, params *models.GetCitiesParams) ([]models.CityDB, int64, error)

	// GetCityByCode получает город по коду и коду страны
	GetCityByCode(ctx context.Context, code, countryCode string) (*models.CityDB, error)

	// CreateCity создает новый город
	CreateCity(ctx context.Context, city *models.CityDB) error

	// UpdateCity обновляет город
	UpdateCity(ctx context.Context, city *models.CityDB) error

	// DeleteCity удаляет город
	DeleteCity(ctx context.Context, code, countryCode string) error
}

// cityRepository реализация репозитория городов
type cityRepository struct {
	db *gorm.DB
}

// CreateCityRepository создает репозиторий городов
func CreateCityRepository(db *gorm.DB) CityRepository {
	return &cityRepository{db: db}
}

// GetCitiesByCountryCode получает города по коду страны
func (r *cityRepository) GetCitiesByCountryCode(ctx context.Context, countryCode string) ([]models.CityDB, error) {
	var cities []models.CityDB
	if err := r.db.WithContext(ctx).Where("country_code = ?", countryCode).Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

// GetCitiesWithFilters получает города с фильтрацией, сортировкой и пагинацией
func (r *cityRepository) GetCitiesWithFilters(ctx context.Context, params *models.GetCitiesParams) ([]models.CityDB, int64, error) {
	var cities []models.CityDB
	var total int64

	query := r.db.WithContext(ctx).Model(&models.CityDB{})

	// Применяем фильтры
	if params.CountryCode != "" {
		query = query.Where("country_code = ?", params.CountryCode)
	}

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
		"code":         true,
		"country_code": true,
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
	type CityWithTotal struct {
		models.CityDB
		Total int64 `gorm:"column:total"`
	}

	var results []CityWithTotal
	if err := query.
		Select("id, code, country_code, created_at, updated_at, COUNT(*) OVER() as total").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(params.Limit).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	// Извлекаем cities и total из результатов
	cities = make([]models.CityDB, len(results))
	if len(results) > 0 {
		total = results[0].Total
		for i, r := range results {
			cities[i] = r.CityDB
		}
	} else {
		total = 0
	}

	return cities, total, nil
}

// GetCityByCode получает город по коду и коду страны
func (r *cityRepository) GetCityByCode(ctx context.Context, code, countryCode string) (*models.CityDB, error) {
	var city models.CityDB
	if err := r.db.WithContext(ctx).Where("code = ? AND country_code = ?", code, countryCode).First(&city).Error; err != nil {
		return nil, err
	}
	return &city, nil
}

// CreateCity создает новый город
func (r *cityRepository) CreateCity(ctx context.Context, city *models.CityDB) error {
	return r.db.WithContext(ctx).Create(city).Error
}

// UpdateCity обновляет город
func (r *cityRepository) UpdateCity(ctx context.Context, city *models.CityDB) error {
	updates := make(map[string]interface{})
	
	// Обновляем только изменяемые поля (исключаем ID, CreatedAt, DeletedAt)
	// Примечание: Code и CountryCode обычно не меняются (они в unique index)
	// Но если нужно обновить, то обновляем их
	if city.Code != "" {
		updates["code"] = city.Code
	}
	if city.CountryCode != "" {
		updates["country_code"] = city.CountryCode
	}
	
	return r.db.WithContext(ctx).Model(city).Updates(updates).Error
}

// DeleteCity удаляет город
func (r *cityRepository) DeleteCity(ctx context.Context, code, countryCode string) error {
	return r.db.WithContext(ctx).Where("code = ? AND country_code = ?", code, countryCode).Delete(&models.CityDB{}).Error
}
