package dictionary

import (
	"dictionary-service/internal/models"
)

// LanguageDictionary содержит все языки с переводами
var LanguageDictionary = map[string]map[models.LanguageCode]string{
	"ru-RU": {
		models.LanguageRU: "Русский",
		models.LanguageKK: "Орыс тілі",
		models.LanguageEN: "Russian",
	},
	"kk-KZ": {
		models.LanguageRU: "Казахский",
		models.LanguageKK: "Қазақ тілі",
		models.LanguageEN: "Kazakh",
	},
	"en-US": {
		models.LanguageRU: "Английский",
		models.LanguageKK: "Ағылшын тілі",
		models.LanguageEN: "English",
	},
}

// GetAllLanguages возвращает список всех языков с переводами для указанного языка
func GetAllLanguages(language models.LanguageCode) []models.Language {
	languages := []models.Language{
		{
			Code: "ru-RU",
			Name: getLanguageName("ru-RU", language),
		},
		{
			Code: "kk-KZ",
			Name: getLanguageName("kk-KZ", language),
		},
		{
			Code: "en-US",
			Name: getLanguageName("en-US", language),
		},
	}
	return languages
}

// getLanguageName возвращает название языка на указанном языке
func getLanguageName(langCode string, targetLanguage models.LanguageCode) string {
	if nameMap, exists := LanguageDictionary[langCode]; exists {
		if name, exists := nameMap[targetLanguage]; exists {
			return name
		}
		// Fallback to English if language not found
		if name, exists := nameMap[models.LanguageEN]; exists {
			return name
		}
	}
	// Fallback to language code itself
	return langCode
}
