package dictionary

import (
	"dictionary-service/internal/models"
)

// MessageCode код системного сообщения
type MessageCode string

const (
	MessageCountriesListReceived MessageCode = "countries.list.received"
	MessageCitiesListReceived    MessageCode = "cities.list.received"
	MessageErrorReceived         MessageCode = "error.received"
	MessageTranslationsReceived  MessageCode = "translations.received"
	MessageErrorsListReceived    MessageCode = "errors.list.received"
	MessageLanguagesListReceived MessageCode = "languages.list.received"
)

// MessageDictionary содержит все системные сообщения на всех языках
var MessageDictionary = map[MessageCode]map[models.LanguageCode]string{
	MessageCountriesListReceived: {
		models.LanguageRU: "Список стран получен",
		models.LanguageKK: "Елдер тізімі алынды",
		models.LanguageEN: "Countries list received",
	},
	MessageCitiesListReceived: {
		models.LanguageRU: "Список городов получен",
		models.LanguageKK: "Қалалар тізімі алынды",
		models.LanguageEN: "Cities list received",
	},
	MessageErrorReceived: {
		models.LanguageRU: "Локализованная ошибка получена",
		models.LanguageKK: "Локализацияланған қате алынды",
		models.LanguageEN: "Localized error received",
	},
	MessageTranslationsReceived: {
		models.LanguageRU: "Переводы получены",
		models.LanguageKK: "Аудармалар алынды",
		models.LanguageEN: "Translations received",
	},
	MessageErrorsListReceived: {
		models.LanguageRU: "Список ошибок получен",
		models.LanguageKK: "Қателер тізімі алынды",
		models.LanguageEN: "Errors list received",
	},
	MessageLanguagesListReceived: {
		models.LanguageRU: "Список языков получен",
		models.LanguageKK: "Тілдер тізімі алынды",
		models.LanguageEN: "Languages list received",
	},
}

// GetMessage возвращает локализованное системное сообщение
func GetMessage(code MessageCode, language models.LanguageCode) string {
	if messageMap, exists := MessageDictionary[code]; exists {
		if message, exists := messageMap[language]; exists {
			return message
		}
		// Fallback to English if language not found
		if message, exists := messageMap[models.LanguageEN]; exists {
			return message
		}
	}
	// Fallback to English default message
	return "Message not found"
}









