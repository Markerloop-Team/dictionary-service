package models

import "fmt"

// LanguageCode представляет поддерживаемые языки
type LanguageCode string

const (
	LanguageRU LanguageCode = "ru"
	LanguageKK LanguageCode = "kk"
	LanguageEN LanguageCode = "en"
)

// DictionaryEntry представляет запись в словаре
type DictionaryEntry struct {
	Key    string            `json:"key"`
	Values map[string]string `json:"values"` // language -> translation
}

// ErrorCode представляет код ошибки
type ErrorCode string

const (
	// Общие ошибки
	ErrorInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorBadRequest     ErrorCode = "BAD_REQUEST"
	ErrorUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrorForbidden      ErrorCode = "FORBIDDEN"
	ErrorNotFound       ErrorCode = "NOT_FOUND"
	ErrorConflict       ErrorCode = "CONFLICT"
	ErrorValidation     ErrorCode = "VALIDATION_ERROR"

	// Auth ошибки
	ErrorUserAlreadyExists       ErrorCode = "USER_ALREADY_EXISTS"
	ErrorUserNotFound            ErrorCode = "USER_NOT_FOUND"
	ErrorInvalidCredentials      ErrorCode = "INVALID_CREDENTIALS"
	ErrorEmailNotVerified        ErrorCode = "EMAIL_NOT_VERIFIED"
	ErrorInvalidToken            ErrorCode = "INVALID_TOKEN"
	ErrorTokenExpired            ErrorCode = "TOKEN_EXPIRED"
	ErrorInvalidVerificationCode ErrorCode = "INVALID_VERIFICATION_CODE"
)

// LocalizedError представляет локализованную ошибку
type LocalizedError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// GetErrorRequest запрос на получение ошибки
// Язык определяется из заголовка Accept-Language
type GetErrorRequest struct {
	Code ErrorCode `json:"code" binding:"required,min=3,max=100"`
	App  *string   `json:"app,omitempty"` // фильтр по приложению (если не указан - общие переводы)
}

// Validate проверяет валидность запроса
func (r *GetErrorRequest) Validate() error {
	if r.Code == "" {
		return fmt.Errorf("код ошибки обязателен")
	}
	if len(r.Code) < 3 || len(r.Code) > 100 {
		return fmt.Errorf("код ошибки должен быть от 3 до 100 символов")
	}
	return nil
}

// GetErrorResponse ответ с локализованной ошибкой
type GetErrorResponse struct {
	Success bool           `json:"success"`
	Data    LocalizedError `json:"data"`
}

// GetTranslationsRequest запрос на получение переводов
// Язык определяется из заголовка Accept-Language
type GetTranslationsRequest struct {
	Codes []string `json:"codes,omitempty"` // если пустой - возвращает все
	Keys  []string `json:"keys,omitempty"`  // алиас для codes (совместимость с go-sdk)
	App   *string  `json:"app,omitempty"`   // фильтр по приложению (если не указан - общие переводы)
}

// MergeCodes объединяет codes и keys в одно поле
func (r *GetTranslationsRequest) MergeCodes() {
	if len(r.Codes) == 0 && len(r.Keys) > 0 {
		r.Codes = r.Keys
	}
}

// Validate проверяет валидность запроса
func (r *GetTranslationsRequest) Validate() error {
	r.MergeCodes()
	for i, code := range r.Codes {
		if code == "" {
			return fmt.Errorf("код перевода не может быть пустым (индекс %d)", i)
		}
		if len(code) > 255 {
			return fmt.Errorf("код перевода слишком длинный (максимум 255 символов): %s", code)
		}
		// Проверка на допустимые символы (буквы, цифры, точки, подчеркивания, дефисы)
		for _, char := range code {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '.' || char == '_' || char == '-') {
				return fmt.Errorf("код перевода содержит недопустимые символы: %s", code)
			}
		}
	}
	return nil
}

// GetTranslationsResponse ответ с переводами
type GetTranslationsResponse struct {
	Success bool                       `json:"success"`
	Data    map[string]DictionaryEntry `json:"data"`
}

// Country представляет страну с переводами
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"` // Название на запрошенном языке
}

// City представляет город с переводами
type City struct {
	Code        string `json:"code"`
	CountryCode string `json:"country_code"`
	Name        string `json:"name"` // Название на запрошенном языке
}

// Language представляет язык
type Language struct {
	Code string `json:"code"` // ru-RU, kk-KZ, en-US
	Name string `json:"name"` // Название на запрошенном языке
}

// GetCountriesRequest запрос на получение стран
type GetCountriesRequest struct {
	// Язык определяется из заголовка Accept-Language
}

// GetCitiesRequest запрос на получение городов
type GetCitiesRequest struct {
	CountryCode string `json:"country_code" binding:"required,min=2,max=3"` // Код страны (обязательное поле)
	// Язык определяется из заголовка Accept-Language
}

// GetCountriesParams параметры для получения стран с фильтрацией, сортировкой и пагинацией
type GetCountriesParams struct {
	Page      int    `json:"page"`       // Номер страницы (по умолчанию: 1) - используется если Cursor не указан
	Limit     int    `json:"limit"`      // Количество записей на странице (по умолчанию: 10, максимум: 100)
	Cursor    string `json:"cursor"`     // ID последнего элемента предыдущей страницы (для cursor-based pagination)
	Search    string `json:"search"`     // Поиск по code, name
	SortBy    string `json:"sort_by"`    // Поле для сортировки (code, name)
	SortOrder string `json:"sort_order"` // Порядок сортировки (asc, desc, по умолчанию: asc)
}

// CountriesResponse ответ со списком стран
type CountriesResponse struct {
	Countries []Country `json:"countries"`
	Count     int       `json:"count"`  // Количество на текущей странице
	Total     int       `json:"total"`  // Общее количество
	Page      int       `json:"page"`   // Текущая страница
	Limit     int       `json:"limit"`  // Лимит на странице
}

// GetCitiesParams параметры для получения городов с фильтрацией, сортировкой и пагинацией
type GetCitiesParams struct {
	Page        int    `json:"page"`         // Номер страницы (по умолчанию: 1) - используется если Cursor не указан
	Limit       int    `json:"limit"`        // Количество записей на странице (по умолчанию: 10, максимум: 100)
	Cursor      string `json:"cursor"`       // ID последнего элемента предыдущей страницы (для cursor-based pagination)
	CountryCode string `json:"country_code"` // Фильтр по коду страны (обязательное поле)
	Search      string `json:"search"`       // Поиск по code, name
	SortBy      string `json:"sort_by"`      // Поле для сортировки (code, name, country_code)
	SortOrder   string `json:"sort_order"`   // Порядок сортировки (asc, desc, по умолчанию: asc)
}

// CitiesResponse ответ со списком городов
type CitiesResponse struct {
	Cities []City `json:"cities"`
	Count  int    `json:"count"`  // Количество на текущей странице
	Total  int    `json:"total"`  // Общее количество
	Page   int    `json:"page"`   // Текущая страница
	Limit  int    `json:"limit"`  // Лимит на странице
}
