package models

import (
	"time"

	"gorm.io/gorm"
)

// TranslationDB универсальная модель для хранения переводов в БД
type TranslationDB struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"uniqueIndex:idx_code_lang_app;size:255;not null" json:"code"` // Ключ/код перевода
	Language  string         `gorm:"uniqueIndex:idx_code_lang_app;size:10;not null" json:"language"`
	Value     string         `gorm:"type:text;not null" json:"value"`                                   // Значение перевода
	App       *string        `gorm:"uniqueIndex:idx_code_lang_app;size:100;index" json:"app,omitempty"` // Приложение (NULL = общий для всех)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName возвращает имя таблицы
func (TranslationDB) TableName() string {
	return "translations"
}

// CountryDB модель для хранения стран в БД
// Переводы хранятся в таблице translations с кодом вида "country.KZ"
type CountryDB struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"uniqueIndex;size:3;not null" json:"code"` // ISO код страны (например, KZ, RU, US)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName возвращает имя таблицы
func (CountryDB) TableName() string {
	return "countries"
}

// CityDB модель для хранения городов в БД
// Переводы хранятся в таблице translations с кодом вида "city.ALA.KZ"
type CityDB struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"uniqueIndex:idx_city_country;size:50;not null" json:"code"`              // Код города
	CountryCode string         `gorm:"uniqueIndex:idx_city_country;size:3;not null;index" json:"country_code"` // Код страны
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName возвращает имя таблицы
func (CityDB) TableName() string {
	return "cities"
}
