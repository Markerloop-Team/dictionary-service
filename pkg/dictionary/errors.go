package dictionary

import (
	"dictionary-service/internal/models"
)

// ErrorDictionary содержит все ошибки на всех языках
var ErrorDictionary = map[models.ErrorCode]map[models.LanguageCode]models.LocalizedError{
	// Общие ошибки
	models.ErrorInternalServer: {
		models.LanguageRU: {
			Code:    models.ErrorInternalServer,
			Message: "Внутренняя ошибка сервера",
			Details: "Произошла непредвиденная ошибка. Попробуйте позже.",
		},
		models.LanguageKK: {
			Code:    models.ErrorInternalServer,
			Message: "Сервердің ішкі қатесі",
			Details: "Күтпеген қате орын алды. Кейінірек көріңіз.",
		},
		models.LanguageEN: {
			Code:    models.ErrorInternalServer,
			Message: "Internal server error",
			Details: "An unexpected error occurred. Please try again later.",
		},
	},

	models.ErrorBadRequest: {
		models.LanguageRU: {
			Code:    models.ErrorBadRequest,
			Message: "Неверный запрос",
			Details: "Проверьте правильность отправленных данных.",
		},
		models.LanguageKK: {
			Code:    models.ErrorBadRequest,
			Message: "Дұрыс емес сұрау",
			Details: "Жіберілген деректердің дұрыстығын тексеріңіз.",
		},
		models.LanguageEN: {
			Code:    models.ErrorBadRequest,
			Message: "Bad request",
			Details: "Please check the data you sent.",
		},
	},

	models.ErrorUnauthorized: {
		models.LanguageRU: {
			Code:    models.ErrorUnauthorized,
			Message: "Не авторизован",
			Details: "Для доступа к этому ресурсу требуется авторизация.",
		},
		models.LanguageKK: {
			Code:    models.ErrorUnauthorized,
			Message: "Авторизацияланбаған",
			Details: "Бұл ресурсқа кіру үшін авторизация қажет.",
		},
		models.LanguageEN: {
			Code:    models.ErrorUnauthorized,
			Message: "Unauthorized",
			Details: "Authorization is required to access this resource.",
		},
	},

	models.ErrorNotFound: {
		models.LanguageRU: {
			Code:    models.ErrorNotFound,
			Message: "Не найдено",
			Details: "Запрашиваемый ресурс не найден.",
		},
		models.LanguageKK: {
			Code:    models.ErrorNotFound,
			Message: "Табылмады",
			Details: "Сұралған ресурс табылмады.",
		},
		models.LanguageEN: {
			Code:    models.ErrorNotFound,
			Message: "Not found",
			Details: "The requested resource was not found.",
		},
	},

	models.ErrorForbidden: {
		models.LanguageRU: {
			Code:    models.ErrorForbidden,
			Message: "Доступ запрещён",
			Details: "У вас нет прав на выполнение этого действия.",
		},
		models.LanguageKK: {
			Code:    models.ErrorForbidden,
			Message: "Қатынауға тыйым салынған",
			Details: "Бұл әрекетті орындауға құқығыңыз жоқ.",
		},
		models.LanguageEN: {
			Code:    models.ErrorForbidden,
			Message: "Access denied",
			Details: "You do not have permission to perform this action.",
		},
	},

	models.ErrorConflict: {
		models.LanguageRU: {
			Code:    models.ErrorConflict,
			Message: "Конфликт",
			Details: "Запрашиваемое действие конфликтует с текущим состоянием ресурса.",
		},
		models.LanguageKK: {
			Code:    models.ErrorConflict,
			Message: "Қайшылық",
			Details: "Сұралған әрекет ресурстың ағымдағы күйімен қайшы келеді.",
		},
		models.LanguageEN: {
			Code:    models.ErrorConflict,
			Message: "Conflict",
			Details: "The requested action conflicts with the current state of the resource.",
		},
	},

	models.ErrorValidation: {
		models.LanguageRU: {
			Code:    models.ErrorValidation,
			Message: "Ошибка валидации",
			Details: "Проверьте правильность заполнения полей.",
		},
		models.LanguageKK: {
			Code:    models.ErrorValidation,
			Message: "Тексеру қатесі",
			Details: "Өрістердің дұрыс толтырылғанын тексеріңіз.",
		},
		models.LanguageEN: {
			Code:    models.ErrorValidation,
			Message: "Validation error",
			Details: "Please check that the fields are filled in correctly.",
		},
	},

	// Auth ошибки
	models.ErrorUserAlreadyExists: {
		models.LanguageRU: {
			Code:    models.ErrorUserAlreadyExists,
			Message: "Пользователь уже существует",
			Details: "Пользователь с таким email уже зарегистрирован в системе.",
		},
		models.LanguageKK: {
			Code:    models.ErrorUserAlreadyExists,
			Message: "Пайдаланушы бұрыннан бар",
			Details: "Осы email-мен пайдаланушы жүйеде бұрыннан тіркелген.",
		},
		models.LanguageEN: {
			Code:    models.ErrorUserAlreadyExists,
			Message: "User already exists",
			Details: "A user with this email is already registered in the system.",
		},
	},

	models.ErrorUserNotFound: {
		models.LanguageRU: {
			Code:    models.ErrorUserNotFound,
			Message: "Пользователь не найден",
			Details: "Пользователь с указанными данными не найден в системе.",
		},
		models.LanguageKK: {
			Code:    models.ErrorUserNotFound,
			Message: "Пайдаланушы табылмады",
			Details: "Көрсетілген деректермен пайдаланушы жүйеде табылмады.",
		},
		models.LanguageEN: {
			Code:    models.ErrorUserNotFound,
			Message: "User not found",
			Details: "User with the specified data was not found in the system.",
		},
	},

	models.ErrorInvalidCredentials: {
		models.LanguageRU: {
			Code:    models.ErrorInvalidCredentials,
			Message: "Неверные учетные данные",
			Details: "Email или пароль указаны неверно. Проверьте правильность ввода.",
		},
		models.LanguageKK: {
			Code:    models.ErrorInvalidCredentials,
			Message: "Дұрыс емес деректер",
			Details: "Email немесе құпия сөз дұрыс көрсетілмеген. Енгізген деректерді тексеріңіз.",
		},
		models.LanguageEN: {
			Code:    models.ErrorInvalidCredentials,
			Message: "Invalid credentials",
			Details: "Email or password is incorrect. Please check your input.",
		},
	},

	models.ErrorEmailNotVerified: {
		models.LanguageRU: {
			Code:    models.ErrorEmailNotVerified,
			Message: "Email не подтвержден",
			Details: "Для входа необходимо подтвердить email. Проверьте почту и перейдите по ссылке.",
		},
		models.LanguageKK: {
			Code:    models.ErrorEmailNotVerified,
			Message: "Email расталмаған",
			Details: "Кіру үшін email-ді растау керек. Поштаны тексеріп, сілтемеге өтіңіз.",
		},
		models.LanguageEN: {
			Code:    models.ErrorEmailNotVerified,
			Message: "Email not verified",
			Details: "You need to verify your email to sign in. Check your email and follow the link.",
		},
	},

	models.ErrorInvalidToken: {
		models.LanguageRU: {
			Code:    models.ErrorInvalidToken,
			Message: "Неверный токен",
			Details: "Предоставленный токен недействителен или поврежден.",
		},
		models.LanguageKK: {
			Code:    models.ErrorInvalidToken,
			Message: "Дұрыс емес токен",
			Details: "Берілген токен жарамсыз немесе бұзылған.",
		},
		models.LanguageEN: {
			Code:    models.ErrorInvalidToken,
			Message: "Invalid token",
			Details: "The provided token is invalid or corrupted.",
		},
	},

	models.ErrorTokenExpired: {
		models.LanguageRU: {
			Code:    models.ErrorTokenExpired,
			Message: "Токен истек",
			Details: "Срок действия токена истек. Войдите в систему заново.",
		},
		models.LanguageKK: {
			Code:    models.ErrorTokenExpired,
			Message: "Токен мерзімі өтті",
			Details: "Токен мерзімі өтті. Жүйеге қайта кіріңіз.",
		},
		models.LanguageEN: {
			Code:    models.ErrorTokenExpired,
			Message: "Token expired",
			Details: "The token has expired. Please sign in again.",
		},
	},

	models.ErrorInvalidVerificationCode: {
		models.LanguageRU: {
			Code:    models.ErrorInvalidVerificationCode,
			Message: "Неверный код подтверждения",
			Details: "Введенный код подтверждения неверен или истек. Запросите новый код.",
		},
		models.LanguageKK: {
			Code:    models.ErrorInvalidVerificationCode,
			Message: "Дұрыс емес растау коды",
			Details: "Енгізілген растау коды дұрыс емес немесе мерзімі өтті. Жаңа код сұраңыз.",
		},
		models.LanguageEN: {
			Code:    models.ErrorInvalidVerificationCode,
			Message: "Invalid verification code",
			Details: "The verification code entered is incorrect or expired. Please request a new code.",
		},
	},
}

// GetError возвращает локализованную ошибку
func GetError(code models.ErrorCode, language models.LanguageCode) models.LocalizedError {
	if errorMap, exists := ErrorDictionary[code]; exists {
		if localizedError, exists := errorMap[language]; exists {
			return localizedError
		}
		// Fallback to English if language not found
		if localizedError, exists := errorMap[models.LanguageEN]; exists {
			return localizedError
		}
	}

	// Fallback to internal server error
	return models.LocalizedError{
		Code:    models.ErrorInternalServer,
		Message: "Internal server error",
		Details: "An unexpected error occurred",
	}
}
