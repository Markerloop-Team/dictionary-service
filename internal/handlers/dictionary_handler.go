package handlers

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"dictionary-service/internal/models"
	"dictionary-service/internal/services"
	"dictionary-service/pkg/dictionary"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	// Импорт go-sdk
	"gitlab.com/markerloop/go-sdk/pkg/response"
)

// maxQueryStringLen caps free-form string params before they hit the
// SQL layer. Prevents pathological inputs (memory pressure, slow
// LIKE) and gives schemathesis fuzzing a soft ceiling.
const maxQueryStringLen = 200

// sanitiseQueryString strips control characters (Postgres TEXT does
// not accept NULL bytes) and caps length. Pure helper — keeps every
// list endpoint defensively consistent without polluting the service
// layer with input policy.
func sanitiseQueryString(in string) string {
	cleaned := strings.Map(func(r rune) rune {
		if r == 0 {
			return -1
		}
		if r < 0x20 && r != '\t' {
			return -1
		}
		return r
	}, in)
	if len(cleaned) > maxQueryStringLen {
		// Обрезаем по границе руны: срез по байтам может разрубить
		// многобайтовый символ и снова получить невалидный UTF-8 —
		// ровно то, от чего этот хелпер и защищает.
		cut := maxQueryStringLen
		for cut > 0 && !utf8.RuneStart(cleaned[cut]) {
			cut--
		}
		cleaned = cleaned[:cut]
	}
	return cleaned
}

// DictionaryHandler обработчик для работы со словарем
type DictionaryHandler struct {
	dictionaryService services.DictionaryService
	log               *logrus.Logger
}

// CreateDictionaryHandler создает хендлер словаря
func CreateDictionaryHandler(dictionaryService services.DictionaryService, log *logrus.Logger) *DictionaryHandler {
	return &DictionaryHandler{
		dictionaryService: dictionaryService,
		log:               log,
	}
}

// GetError возвращает локализованную ошибку
// @Summary Получить локализованную ошибку
// @Description Возвращает локализованное сообщение об ошибке по коду. Язык из Accept-Language.
// @Tags dictionary
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Язык (ru, kk, en)" default(en)
// @Param request body models.GetErrorRequest true "Код ошибки"
// @Success 200 {object} models.LocalizedErrorSuccessResponse
// @Failure 400 {object} models.ErrorResponseSwagger
// @Router /dictionary/error [post]
func (h *DictionaryHandler) GetError(c *gin.Context) {
	var req models.GetErrorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Ошибка валидации запроса:", err) // 4xx по вине вызывающего — не Sentry (SDL-215)
		response.BadRequest(c, "Неверный формат запроса")
		return
	}

	// Дополнительная валидация
	if err := req.Validate(); err != nil {
		h.log.Warnf("Ошибка валидации: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	// Получаем язык из заголовка Accept-Language
	language := h.getLanguageFromHeader(c)

	// Получаем локализованную ошибку с учетом приложения
	localizedError := h.dictionaryService.GetError(req.Code, language, req.App)

	// Получаем переведенное сообщение
	message := dictionary.GetMessage(dictionary.MessageErrorReceived, language)
	response.Success(c, localizedError, message)
}

// GetTranslations возвращает переводы
// @Summary Получить переводы
// @Description Возвращает переводы по списку кодов. Язык из Accept-Language.
// @Tags dictionary
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Язык (ru, kk, en)" default(en)
// @Param request body models.GetTranslationsRequest true "Коды переводов"
// @Success 200 {object} models.TranslationsSuccessResponse
// @Failure 400 {object} models.ErrorResponseSwagger
// @Router /dictionary/translations [post]
func (h *DictionaryHandler) GetTranslations(c *gin.Context) {
	var req models.GetTranslationsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Ошибка валидации запроса:", err) // 4xx по вине вызывающего — не Sentry (SDL-215)
		response.BadRequest(c, "Неверный формат запроса")
		return
	}

	// Дополнительная валидация
	if err := req.Validate(); err != nil {
		h.log.Warnf("Ошибка валидации: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	// Получаем язык из заголовка Accept-Language
	language := h.getLanguageFromHeader(c)

	// Получаем переводы с учетом приложения
	translations := h.dictionaryService.GetTranslations(language, req.Codes, req.App)

	// Получаем переведенное сообщение
	message := dictionary.GetMessage(dictionary.MessageTranslationsReceived, language)
	response.Success(c, translations, message)
}

// getLanguageFromHeader извлекает и нормализует язык из заголовка Accept-Language
// Поддерживает форматы: "ru", "ru-RU", "kk", "kk-KZ", "en", "en-US", "en-GB"
// Стандартизирует к форматам: ru-RU, kk-KZ, en-US
func (h *DictionaryHandler) getLanguageFromHeader(c *gin.Context) models.LanguageCode {
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		h.log.Debug("Заголовок Accept-Language не указан, используется английский по умолчанию")
		return models.LanguageEN
	}

	// Нормализуем язык: извлекаем базовый код языка (первые 2 символа)
	lang := acceptLanguage
	if len(acceptLanguage) >= 2 {
		lang = acceptLanguage[:2]
	}

	// Определяем язык по базовому коду
	switch lang {
	case "kk":
		return models.LanguageKK
	case "ru":
		return models.LanguageRU
	case "en":
		return models.LanguageEN
	default:
		// По умолчанию английский для несуществующих языков
		h.log.Warnf("Неизвестный язык из заголовка Accept-Language: %s, используется английский по умолчанию", acceptLanguage)
		return models.LanguageEN
	}
}

// GetCountries возвращает список стран с поддержкой фильтрации, сортировки и пагинации
// @Summary Получить список стран
// @Description Возвращает страны с пагинацией, фильтрацией и сортировкой. Язык из Accept-Language.
// @Tags dictionary
// @Produce json
// @Param Accept-Language header string false "Язык (ru, kk, en)" default(en)
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Записей на странице (макс 100)" default(10)
// @Param search query string false "Поиск по code, name"
// @Param sort_by query string false "Поле сортировки (code, name)" default(code)
// @Param sort_order query string false "Порядок (asc, desc)" default(asc)
// @Success 200 {object} models.CountriesSuccessResponse
// @Failure 500 {object} models.ErrorResponseSwagger
// @Router /dictionary/countries [get]
func (h *DictionaryHandler) GetCountries(c *gin.Context) {
	// Получаем язык из заголовка Accept-Language
	language := h.getLanguageFromHeader(c)

	// Парсинг query параметров
	params := &models.GetCountriesParams{}

	// Парсинг page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			params.Page = 1
		}
	} else {
		params.Page = 1
	}

	// Парсинг limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100
			}
			params.Limit = limit
		} else {
			params.Limit = 10
		}
	} else {
		params.Limit = 10
	}

	// Парсинг search — sanitise to keep NULL bytes / control chars
	// out of Postgres LIKE patterns.
	params.Search = sanitiseQueryString(c.Query("search"))

	// Парсинг sort_by — sanitise + repository allow-lists to "code".
	params.SortBy = sanitiseQueryString(c.DefaultQuery("sort_by", "code"))

	// Парсинг sort_order — sanitise + repository normalises asc/desc.
	params.SortOrder = sanitiseQueryString(c.DefaultQuery("sort_order", "asc"))

	// Вызов сервиса с фильтрами
	countries, err := h.dictionaryService.GetCountriesWithFilters(language, params)
	if err != nil {
		h.log.Errorf("Ошибка получения стран: %v", err)
		response.InternalServerError(c, "Ошибка получения списка стран")
		return
	}

	// Вычисляем total_pages
	totalPages := (countries.Total + countries.Limit - 1) / countries.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	// Формируем метаданные для пагинации
	meta := &response.Meta{
		Page:       countries.Page,
		Limit:      countries.Limit,
		Total:      countries.Total,
		TotalPages: totalPages,
	}

	// Получаем переведенное сообщение
	message := dictionary.GetMessage(dictionary.MessageCountriesListReceived, language)
	response.SuccessWithMeta(c, countries, message, meta)
}

// GetCities возвращает список городов с поддержкой фильтрации, сортировки и пагинации
// @Summary Получить список городов
// @Description Возвращает города по коду страны с пагинацией и фильтрацией. Язык из Accept-Language.
// @Tags dictionary
// @Produce json
// @Param Accept-Language header string false "Язык (ru, kk, en)" default(en)
// @Param country_code query string true "Код страны (обязательно)"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Записей на странице (макс 100)" default(10)
// @Param search query string false "Поиск по code, name"
// @Param sort_by query string false "Поле сортировки (code, name)" default(code)
// @Param sort_order query string false "Порядок (asc, desc)" default(asc)
// @Success 200 {object} models.CitiesSuccessResponse
// @Failure 400 {object} models.ErrorResponseSwagger
// @Failure 500 {object} models.ErrorResponseSwagger
// @Router /dictionary/cities [get]
func (h *DictionaryHandler) GetCities(c *gin.Context) {
	// Получаем язык из заголовка Accept-Language
	language := h.getLanguageFromHeader(c)

	// Парсинг query параметров
	params := &models.GetCitiesParams{}

	// Парсинг country_code (обязательное поле)
	// SDL-315: country_code так же уходит в запрос к Postgres, как search
	// и sort_*, поэтому санитайзится наравне с ними. NULL-байт или
	// невалидный UTF-8 здесь роняли выборку в 500. Пропущен был потому,
	// что параметр есть только у /cities, а /cities был недостижим через
	// шлюз (метод разошёлся) — фаззер до него не доходил.
	params.CountryCode = sanitiseQueryString(c.Query("country_code"))
	if params.CountryCode == "" {
		response.BadRequest(c, "country_code обязателен")
		return
	}

	// Парсинг page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			params.Page = 1
		}
	} else {
		params.Page = 1
	}

	// Парсинг limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100
			}
			params.Limit = limit
		} else {
			params.Limit = 10
		}
	} else {
		params.Limit = 10
	}

	// Парсинг search — sanitise to keep NULL bytes / control chars
	// out of Postgres LIKE patterns.
	params.Search = sanitiseQueryString(c.Query("search"))

	// Парсинг sort_by
	params.SortBy = sanitiseQueryString(c.DefaultQuery("sort_by", "code"))

	// Парсинг sort_order
	params.SortOrder = sanitiseQueryString(c.DefaultQuery("sort_order", "asc"))

	// Вызов сервиса с фильтрами
	cities, err := h.dictionaryService.GetCitiesWithFilters(language, params)
	if err != nil {
		h.log.Errorf("Ошибка получения городов: %v", err)
		response.InternalServerError(c, "Ошибка получения списка городов")
		return
	}

	// Вычисляем total_pages
	totalPages := (cities.Total + cities.Limit - 1) / cities.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	// Формируем метаданные для пагинации
	meta := &response.Meta{
		Page:       cities.Page,
		Limit:      cities.Limit,
		Total:      cities.Total,
		TotalPages: totalPages,
	}

	// Получаем переведенное сообщение
	message := dictionary.GetMessage(dictionary.MessageCitiesListReceived, language)
	response.SuccessWithMeta(c, cities, message, meta)
}

// GetAllErrors возвращает список всех ошибок
// @Summary Получить все ошибки
// @Description Возвращает список всех локализованных ошибок. Язык из Accept-Language.
// @Tags dictionary
// @Produce json
// @Param Accept-Language header string false "Язык (ru, kk, en)" default(en)
// @Success 200 {object} models.ErrorsListSuccessResponse
// @Router /dictionary/errors [get]
func (h *DictionaryHandler) GetAllErrors(c *gin.Context) {
	// Получаем язык из заголовка Accept-Language
	language := h.getLanguageFromHeader(c)

	// Получаем все ошибки
	errors := h.dictionaryService.GetAllErrors(language)

	// Получаем переведенное сообщение
	message := dictionary.GetMessage(dictionary.MessageErrorsListReceived, language)
	response.Success(c, errors, message)
}

// GetAllLanguages возвращает список всех языков
// @Summary Получить все языки
// @Description Возвращает список всех поддерживаемых языков. Язык из Accept-Language.
// @Tags dictionary
// @Produce json
// @Param Accept-Language header string false "Язык (ru, kk, en)" default(en)
// @Success 200 {object} models.LanguagesSuccessResponse
// @Router /dictionary/languages [get]
func (h *DictionaryHandler) GetAllLanguages(c *gin.Context) {
	// Получаем язык из заголовка Accept-Language
	language := h.getLanguageFromHeader(c)

	// Получаем все языки
	languages := h.dictionaryService.GetAllLanguages(language)

	// Получаем переведенное сообщение
	message := dictionary.GetMessage(dictionary.MessageLanguagesListReceived, language)
	response.Success(c, languages, message)
}

// HealthCheck проверка здоровья сервиса
func (h *DictionaryHandler) HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"service": "dictionary-service",
		"status":  "ok",
		"version": "1.0.0",
	}, "Сервис работает")
}
