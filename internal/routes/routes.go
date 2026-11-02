package routes

import (
	"dictionary-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Импорт go-sdk
	sdkMiddleware "gitlab.com/markerloop/go-sdk/pkg/middleware"
	"gitlab.com/markerloop/go-sdk/pkg/requestid"

	// Swagger docs
	_ "dictionary-service/docs"
)

// SetupRoutes настраивает все роуты приложения
func SetupRoutes(r *gin.Engine, dictionaryHandler *handlers.DictionaryHandler, log *logrus.Logger) {
	// Middleware из SDK
	r.Use(requestid.Middleware()) // SDL-87: сквозной X-Request-Id
	r.Use(sdkMiddleware.Recovery(log))
	r.Use(sdkMiddleware.Logger(log))
	r.Use(sdkMiddleware.CORS())

	// Health check на корне: gw ProxyHealthCheck всегда проксирует на /health
	// сервиса (SDL-206). /api/health ниже оставлен — на него смотрят k8s-пробы
	// (dev-infra/apps/dictionary-service/base/deployment.yaml).
	r.GET("/health", sdkMiddleware.HealthCheck("dictionary-service", "1.0.0"))

	// Swagger UI (на корневом уровне для совместимости с gateway proxy)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API группа
	api := r.Group("/api")
	{
		// Health check из SDK
		api.GET("/health", sdkMiddleware.HealthCheck("dictionary-service", "1.0.0"))
		// API v1 роуты
		v1 := api.Group("/v1")
		{
			// Dictionary роуты
			dictionary := v1.Group("/dictionary")
			{
				dictionary.POST("/error", dictionaryHandler.GetError)               // Получить локализованную ошибку
				dictionary.GET("/errors", dictionaryHandler.GetAllErrors)           // Получить список всех ошибок
				dictionary.POST("/translations", dictionaryHandler.GetTranslations) // Получить переводы
				dictionary.POST("/words", dictionaryHandler.GetTranslations)        // Алиас для go-sdk (использует /words)
				dictionary.GET("/countries", dictionaryHandler.GetCountries)        // Получить список стран (с пагинацией, фильтрацией, сортировкой)
				dictionary.GET("/cities", dictionaryHandler.GetCities)              // Получить список городов (с пагинацией, фильтрацией, сортировкой)
				dictionary.GET("/languages", dictionaryHandler.GetAllLanguages)     // Получить список всех языков
			}
		}
	}
}
