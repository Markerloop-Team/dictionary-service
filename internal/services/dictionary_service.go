package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"dictionary-service/internal/models"
	"dictionary-service/internal/repository"
	"dictionary-service/pkg/dictionary"

	"github.com/sirupsen/logrus"

	// Импорт go-sdk
	"gitlab.com/markerloop/go-sdk/pkg/redis"
)

// DictionaryService интерфейс для работы со словарем
type DictionaryService interface {
	GetError(code models.ErrorCode, language models.LanguageCode, app *string) models.LocalizedError
	GetTranslations(language models.LanguageCode, codes []string, app *string) map[string]models.DictionaryEntry
	GetCountries(language models.LanguageCode) ([]models.Country, error)
	GetCountriesWithFilters(language models.LanguageCode, params *models.GetCountriesParams) (*models.CountriesResponse, error)
	GetCities(language models.LanguageCode, countryCode string) ([]models.City, error)
	GetCitiesWithFilters(language models.LanguageCode, params *models.GetCitiesParams) (*models.CitiesResponse, error)
	GetAllErrors(language models.LanguageCode) []models.LocalizedError
	GetAllLanguages(language models.LanguageCode) []models.Language
}

// dictionaryService реализация сервиса словаря
type dictionaryService struct {
	translationRepo repository.TranslationRepository
	countryRepo     repository.CountryRepository
	cityRepo        repository.CityRepository
	redisService    *redis.RedisService
	log             *logrus.Logger
	cacheTTL        time.Duration
}

// CreateDictionaryService создает сервис словаря
func CreateDictionaryService(
	translationRepo repository.TranslationRepository,
	countryRepo repository.CountryRepository,
	cityRepo repository.CityRepository,
	redisService *redis.RedisService,
	cacheTTL time.Duration,
	log *logrus.Logger,
) DictionaryService {
	return &dictionaryService{
		translationRepo: translationRepo,
		countryRepo:     countryRepo,
		cityRepo:        cityRepo,
		redisService:    redisService,
		log:             log,
		cacheTTL:        cacheTTL,
	}
}

// GetError возвращает локализованную ошибку
func (s *dictionaryService) GetError(code models.ErrorCode, language models.LanguageCode, app *string) models.LocalizedError {
	ctx := context.Background()
	appStr := ""
	if app != nil {
		appStr = *app
	}
	cacheKey := fmt.Sprintf("error:%s:%s:%s", code, language, appStr)

	// Пытаемся получить из кэша
	if s.redisService != nil {
		cached, err := s.redisService.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var localizedError models.LocalizedError
			if err := json.Unmarshal([]byte(cached), &localizedError); err == nil {
				s.log.Debugf("Получена ошибка из кэша: %s", code)
				return localizedError
			}
		}
	}

	// Получаем из таблицы translations
	translations, err := s.translationRepo.GetTranslations(ctx, []string{string(code)}, string(language), app)
	if err == nil && len(translations) > 0 {
		// Берем первый перевод (уже отсортирован с приоритетом специфичных для app)
		translation := translations[0]
		localizedError := models.LocalizedError{
			Code:    models.ErrorCode(translation.Code),
			Message: translation.Value,
			Details: translation.Value, // Используем Value вместо Details
		}

		// Сохраняем в кэш
		if s.redisService != nil {
			if data, err := json.Marshal(localizedError); err == nil {
				s.redisService.Set(ctx, cacheKey, string(data), s.cacheTTL)
			}
		}

		return localizedError
	}

	if err != nil {
		s.log.Warnf("Ошибка получения перевода из БД для кода: %s, языка: %s, app: %v: %v", code, language, app, err)
	} else {
		s.log.Warnf("Перевод не найден в БД для кода: %s, языка: %s, app: %v", code, language, app)
	}

	// Fallback на статический словарь
	s.log.Debugf("Используется fallback для ошибки: %s", code)
	return dictionary.GetError(code, language)
}

// GetTranslations возвращает переводы для указанного языка и приложения
func (s *dictionaryService) GetTranslations(language models.LanguageCode, codes []string, app *string) map[string]models.DictionaryEntry {
	ctx := context.Background()
	result := make(map[string]models.DictionaryEntry)

	appStr := ""
	if app != nil {
		appStr = *app
	}

	// Формируем ключ кэша
	cacheKey := fmt.Sprintf("dict:%s:%s:%v", language, appStr, codes)
	if len(codes) == 0 {
		cacheKey = fmt.Sprintf("dict:%s:%s:all", language, appStr)
	}

	// Пытаемся получить из кэша
	if s.redisService != nil {
		cached, err := s.redisService.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var cachedDict map[string]models.DictionaryEntry
			if err := json.Unmarshal([]byte(cached), &cachedDict); err == nil {
				s.log.Debugf("Получен словарь из кэша для языка: %s, app: %v", language, app)
				return cachedDict
			}
		}
	}

	// Получаем из таблицы translations
	var translations []models.TranslationDB
	var err error

	if len(codes) == 0 {
		translations, err = s.translationRepo.GetAllTranslations(ctx, string(language), app)
	} else {
		translations, err = s.translationRepo.GetTranslations(ctx, codes, string(language), app)
	}

	if err == nil && len(translations) > 0 {
		// Преобразуем записи БД в формат API
		dictMap := make(map[string]map[string]string) // code -> language -> value

		for _, translation := range translations {
			if dictMap[translation.Code] == nil {
				dictMap[translation.Code] = make(map[string]string)
			}
			dictMap[translation.Code][translation.Language] = translation.Value
		}

		// Преобразуем в формат DictionaryEntry
		for code, values := range dictMap {
			result[code] = models.DictionaryEntry{
				Key:    code,
				Values: values,
			}
		}

		// Сохраняем в кэш
		if s.redisService != nil {
			if data, err := json.Marshal(result); err == nil {
				s.redisService.Set(ctx, cacheKey, string(data), s.cacheTTL)
			}
		}

		s.log.Infof("Получено %d записей словаря из БД для языка: %s, app: %v", len(result), language, app)
		return result
	}

	if err != nil {
		s.log.Warnf("Ошибка получения словаря из БД: %v", err)
	} else {
		s.log.Warnf("Словарь не найден в БД для языка: %s, кодов: %v, app: %v", language, codes, app)
	}

	// Fallback на статический словарь
	s.log.Debugf("Используется fallback словарь для языка: %s", language)
	return s.getFallbackDictionary(language, codes)
}

// getFallbackDictionary возвращает fallback словарь
func (s *dictionaryService) getFallbackDictionary(language models.LanguageCode, keys []string) map[string]models.DictionaryEntry {
	result := make(map[string]models.DictionaryEntry)

	// Примеры записей словаря
	exampleEntries := map[string]models.DictionaryEntry{
		"welcome": {
			Key: "welcome",
			Values: map[string]string{
				string(models.LanguageRU): "Добро пожаловать",
				string(models.LanguageKK): "Қош келдіңіз",
				string(models.LanguageEN): "Welcome",
			},
		},
		"goodbye": {
			Key: "goodbye",
			Values: map[string]string{
				string(models.LanguageRU): "До свидания",
				string(models.LanguageKK): "Сау болыңыз",
				string(models.LanguageEN): "Goodbye",
			},
		},
		"error": {
			Key: "error",
			Values: map[string]string{
				string(models.LanguageRU): "Ошибка",
				string(models.LanguageKK): "Қате",
				string(models.LanguageEN): "Error",
			},
		},
	}

	// Если ключи не указаны, возвращаем все
	if len(keys) == 0 {
		return exampleEntries
	}

	// Возвращаем только запрошенные ключи
	for _, key := range keys {
		if entry, exists := exampleEntries[key]; exists {
			result[key] = entry
		} else {
			s.log.Warnf("Ключ словаря не найден: %s", key)
		}
	}

	return result
}

// GetCountries возвращает список стран с переводами для указанного языка
func (s *dictionaryService) GetCountries(language models.LanguageCode) ([]models.Country, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("countries:%s", language)

	// Пытаемся получить из кэша
	if s.redisService != nil {
		cached, err := s.redisService.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var countries []models.Country
			if err := json.Unmarshal([]byte(cached), &countries); err == nil {
				s.log.Debugf("Получены страны из кэша для языка: %s", language)
				return countries, nil
			}
		}
	}

	// Получаем страны из БД с пагинацией (используем большой лимит для получения всех)
	params := &models.GetCountriesParams{
		Page:  1,
		Limit: 1000, // Большой лимит для получения всех стран
	}
	countriesDB, _, err := s.countryRepo.GetAllCountriesWithFilters(ctx, params)
	if err != nil {
		s.log.Warnf("Ошибка получения стран из БД: %v", err)
		return nil, err
	}

	// Формируем коды для получения переводов: "country.KZ", "country.RU", и т.д.
	codes := make([]string, 0, len(countriesDB))
	for _, countryDB := range countriesDB {
		codes = append(codes, fmt.Sprintf("country.%s", countryDB.Code))
	}

	// Получаем переводы для всех стран на всех языках
	translationsMap, err := s.translationRepo.GetTranslationsByCodes(ctx, codes, nil)
	if err != nil {
		s.log.Warnf("Ошибка получения переводов стран: %v", err)
		// Продолжаем работу, даже если переводы не найдены
		translationsMap = make(map[string]map[string]string)
	}

	// Преобразуем в формат API
	countries := make([]models.Country, 0, len(countriesDB))
	for _, countryDB := range countriesDB {
		code := fmt.Sprintf("country.%s", countryDB.Code)
		names := translationsMap[code]
		if names == nil {
			names = make(map[string]string)
		}

		// Получаем название на запрошенном языке
		name := names[string(language)]
		if name == "" {
			// Fallback на английский, если запрошенный язык отсутствует
			name = names[string(models.LanguageEN)]
			if name == "" {
				// Если перевода нет вообще, используем код страны
				name = countryDB.Code
			}
		}

		country := models.Country{
			Code: countryDB.Code,
			Name: name,
		}
		countries = append(countries, country)
	}

	// Сохраняем в кэш
	if s.redisService != nil {
		if data, err := json.Marshal(countries); err == nil {
			s.redisService.Set(ctx, cacheKey, string(data), s.cacheTTL)
		}
	}

	s.log.Infof("Получено %d стран из БД для языка: %s", len(countries), language)
	return countries, nil
}

// GetCities возвращает список городов для указанной страны с переводами для указанного языка
func (s *dictionaryService) GetCities(language models.LanguageCode, countryCode string) ([]models.City, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cities:%s:%s", countryCode, language)

	// Пытаемся получить из кэша
	if s.redisService != nil {
		cached, err := s.redisService.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var cities []models.City
			if err := json.Unmarshal([]byte(cached), &cities); err == nil {
				s.log.Debugf("Получены города из кэша для страны: %s, языка: %s", countryCode, language)
				return cities, nil
			}
		}
	}

	// Получаем города из БД с пагинацией (используем большой лимит для получения всех)
	params := &models.GetCitiesParams{
		Page:        1,
		Limit:       1000, // Большой лимит для получения всех городов
		CountryCode: countryCode,
	}
	citiesDB, _, err := s.cityRepo.GetCitiesWithFilters(ctx, params)
	if err != nil {
		s.log.Warnf("Ошибка получения городов из БД для страны %s: %v", countryCode, err)
		return nil, err
	}

	// Формируем коды для получения переводов: "city.ALA.KZ", "city.AST.KZ", и т.д.
	codes := make([]string, 0, len(citiesDB))
	for _, cityDB := range citiesDB {
		codes = append(codes, fmt.Sprintf("city.%s.%s", cityDB.Code, cityDB.CountryCode))
	}

	// Получаем переводы для всех городов на всех языках
	translationsMap, err := s.translationRepo.GetTranslationsByCodes(ctx, codes, nil)
	if err != nil {
		s.log.Warnf("Ошибка получения переводов городов: %v", err)
		// Продолжаем работу, даже если переводы не найдены
		translationsMap = make(map[string]map[string]string)
	}

	// Преобразуем в формат API
	cities := make([]models.City, 0, len(citiesDB))
	for _, cityDB := range citiesDB {
		code := fmt.Sprintf("city.%s.%s", cityDB.Code, cityDB.CountryCode)
		names := translationsMap[code]
		if names == nil {
			names = make(map[string]string)
		}

		// Получаем название на запрошенном языке
		name := names[string(language)]
		if name == "" {
			// Fallback на английский, если запрошенный язык отсутствует
			name = names[string(models.LanguageEN)]
			if name == "" {
				// Если перевода нет вообще, используем код города
				name = cityDB.Code
			}
		}

		city := models.City{
			Code:        cityDB.Code,
			CountryCode: cityDB.CountryCode,
			Name:        name,
		}
		cities = append(cities, city)
	}

	// Сохраняем в кэш
	if s.redisService != nil {
		if data, err := json.Marshal(cities); err == nil {
			s.redisService.Set(ctx, cacheKey, string(data), s.cacheTTL)
		}
	}

	s.log.Infof("Получено %d городов из БД для страны: %s, языка: %s", len(cities), countryCode, language)
	return cities, nil
}

// GetCountriesWithFilters возвращает список стран с фильтрацией, сортировкой и пагинацией
func (s *dictionaryService) GetCountriesWithFilters(language models.LanguageCode, params *models.GetCountriesParams) (*models.CountriesResponse, error) {
	ctx := context.Background()

	// Валидация и нормализация параметров
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Получаем страны из БД с фильтрами
	countriesDB, total, err := s.countryRepo.GetAllCountriesWithFilters(ctx, params)
	if err != nil {
		s.log.Warnf("Ошибка получения стран из БД: %v", err)
		return nil, err
	}

	// Формируем коды для получения переводов
	codes := make([]string, 0, len(countriesDB))
	for _, countryDB := range countriesDB {
		codes = append(codes, fmt.Sprintf("country.%s", countryDB.Code))
	}

	// Получаем переводы для отфильтрованных стран
	translationsMap, err := s.translationRepo.GetTranslationsByCodes(ctx, codes, nil)
	if err != nil {
		s.log.Warnf("Ошибка получения переводов стран: %v", err)
		translationsMap = make(map[string]map[string]string)
	}

	// Преобразуем в формат API
	countries := make([]models.Country, 0, len(countriesDB))
	for _, countryDB := range countriesDB {
		code := fmt.Sprintf("country.%s", countryDB.Code)
		names := translationsMap[code]
		if names == nil {
			names = make(map[string]string)
		}

		name := names[string(language)]
		if name == "" {
			name = names[string(models.LanguageEN)]
			if name == "" {
				name = countryDB.Code
			}
		}

		// Применяем фильтр по search для name (если указан)
		if params.Search != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(params.Search)) {
			continue
		}

		country := models.Country{
			Code: countryDB.Code,
			Name: name,
		}
		countries = append(countries, country)
	}

	// Обновляем total с учетом фильтрации по name
	if params.Search != "" {
		total = int64(len(countries))
	}

	return &models.CountriesResponse{
		Countries: countries,
		Count:     len(countries),
		Total:     int(total),
		Page:      params.Page,
		Limit:     params.Limit,
	}, nil
}

// GetCitiesWithFilters возвращает список городов с фильтрацией, сортировкой и пагинацией
func (s *dictionaryService) GetCitiesWithFilters(language models.LanguageCode, params *models.GetCitiesParams) (*models.CitiesResponse, error) {
	ctx := context.Background()

	// Валидация и нормализация параметров
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Получаем города из БД с фильтрами
	citiesDB, total, err := s.cityRepo.GetCitiesWithFilters(ctx, params)
	if err != nil {
		s.log.Warnf("Ошибка получения городов из БД: %v", err)
		return nil, err
	}

	// Формируем коды для получения переводов
	codes := make([]string, 0, len(citiesDB))
	for _, cityDB := range citiesDB {
		codes = append(codes, fmt.Sprintf("city.%s.%s", cityDB.Code, cityDB.CountryCode))
	}

	// Получаем переводы для отфильтрованных городов
	translationsMap, err := s.translationRepo.GetTranslationsByCodes(ctx, codes, nil)
	if err != nil {
		s.log.Warnf("Ошибка получения переводов городов: %v", err)
		translationsMap = make(map[string]map[string]string)
	}

	// Преобразуем в формат API
	cities := make([]models.City, 0, len(citiesDB))
	for _, cityDB := range citiesDB {
		code := fmt.Sprintf("city.%s.%s", cityDB.Code, cityDB.CountryCode)
		names := translationsMap[code]
		if names == nil {
			names = make(map[string]string)
		}

		name := names[string(language)]
		if name == "" {
			name = names[string(models.LanguageEN)]
			if name == "" {
				name = cityDB.Code
			}
		}

		// Применяем фильтр по search для name (если указан)
		if params.Search != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(params.Search)) {
			continue
		}

		city := models.City{
			Code:        cityDB.Code,
			CountryCode: cityDB.CountryCode,
			Name:        name,
		}
		cities = append(cities, city)
	}

	// Обновляем total с учетом фильтрации по name
	if params.Search != "" {
		total = int64(len(cities))
	}

	return &models.CitiesResponse{
		Cities: cities,
		Count:  len(cities),
		Total:  int(total),
		Page:   params.Page,
		Limit:  params.Limit,
	}, nil
}

// GetAllErrors возвращает список всех ошибок с переводами для указанного языка
func (s *dictionaryService) GetAllErrors(language models.LanguageCode) []models.LocalizedError {
	errors := make([]models.LocalizedError, 0)

	// Получаем все коды ошибок из статического словаря
	for errorCode := range dictionary.ErrorDictionary {
		error := s.GetError(errorCode, language, nil)
		errors = append(errors, error)
	}

	return errors
}

// GetAllLanguages возвращает список всех языков с переводами для указанного языка
func (s *dictionaryService) GetAllLanguages(language models.LanguageCode) []models.Language {
	return dictionary.GetAllLanguages(language)
}
