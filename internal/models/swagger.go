package models

// Swagger response wrappers

type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

type ErrorResponseSwagger struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Ошибка"`
	Error   string `json:"error" example:"bad_request"`
}

type LocalizedErrorSuccessResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message"`
	Data    *LocalizedError `json:"data"`
}

type TranslationsSuccessResponse struct {
	Success bool                       `json:"success" example:"true"`
	Message string                     `json:"message"`
	Data    map[string]DictionaryEntry `json:"data"`
}

type CountriesSuccessResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message"`
	Data    *CountriesResponse `json:"data"`
	Meta    *Meta              `json:"meta"`
}

type CitiesSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message"`
	Data    *CitiesResponse `json:"data"`
	Meta    *Meta           `json:"meta"`
}

type ErrorsListSuccessResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message"`
	Data    []LocalizedError `json:"data"`
}

type LanguagesSuccessResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message"`
	Data    []Language `json:"data"`
}
