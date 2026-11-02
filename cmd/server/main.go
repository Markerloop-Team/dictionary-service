package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"dictionary-service/internal/config"
	"dictionary-service/internal/handlers"
	"dictionary-service/internal/migrations"
	"dictionary-service/internal/models"
	"dictionary-service/internal/repository"
	"dictionary-service/internal/routes"
	"dictionary-service/internal/services"

	// Импорт go-sdk
	"gitlab.com/markerloop/go-sdk/pkg/database"
	"gitlab.com/markerloop/go-sdk/pkg/logger"
	sentrysdk "gitlab.com/markerloop/go-sdk/pkg/sentry"
	"gitlab.com/markerloop/go-sdk/pkg/redis"
)

// @title Dictionary Service API
// @version 1.0
// @description Сервис словарей, локализации и справочных данных Markerloop
// @BasePath /api/v1
func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения по умолчанию")
	}

	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем логгер SDK
	loggerInstance := logger.InitLogger("dictionary-service")

	// Sentry (SDL-86): включается только когда контур задал SENTRY_DSN;
	// без переменной — рабочий no-op, сервис стартует как раньше.
	sentrysdk.Init("dictionary-service", loggerInstance)
	defer sentrysdk.Flush()

	// Инициализируем подключение к PostgreSQL с миграциями через SDK
	dbConfig := &database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Name:            cfg.Database.Name,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	db, err := database.NewDatabaseWithMigration(
		dbConfig,
		loggerInstance,
		&models.TranslationDB{}, // Универсальная таблица переводов
		&models.CountryDB{},     // Таблица стран
		&models.CityDB{},        // Таблица городов
	)
	if err != nil {
		loggerInstance.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			loggerInstance.Errorf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	// Инициализируем Redis (опционально, не критично если недоступен)
	var redisService *redis.RedisService
	redisAddr := cfg.Redis.Addr
	if redisAddr != "" {
		// Временно устанавливаем переменные окружения для Redis
		os.Setenv("REDIS_ADDR", cfg.Redis.Addr)
		os.Setenv("REDIS_PASSWORD", cfg.Redis.Password)
		os.Setenv("REDIS_DB", fmt.Sprintf("%d", cfg.Redis.DB))
		redisService = redis.NewRedisService(loggerInstance)
	} else {
		loggerInstance.Warn("Redis не настроен, кэширование отключено")
	}

	// Создаем индексы для оптимизации производительности
	loggerInstance.Info("Создание индексов для оптимизации производительности...")
	if err := migrations.CreatePerformanceIndexes(db.DB); err != nil {
		loggerInstance.Warn("Ошибка создания индексов (может быть уже выполнена):", err)
	} else {
		loggerInstance.Info("Создание индексов завершено")
	}

	// Применяем seed данные (ON CONFLICT — безопасно при повторном запуске)
	loggerInstance.Info("Применение seed данных...")
	if err := migrations.RunSeeds(db.DB, loggerInstance); err != nil {
		loggerInstance.Warn("Ошибка применения seed данных:", err)
	}

	// SDL-313: чистка строк app='web-app'. Идёт ПОСЛЕ сидов намеренно —
	// так инвариант «после старта пода таких строк нет» держится вне
	// зависимости от того, что успел записать сидер.
	if removed, err := migrations.DropWebAppTranslations(db.DB); err != nil {
		loggerInstance.Warn("Ошибка чистки translations app='web-app':", err)
	} else if removed > 0 {
		loggerInstance.Infof("Удалено строк translations app='web-app': %d", removed)
	}

	// Создаем репозитории
	translationRepo := repository.CreateTranslationRepository(db.DB)
	countryRepo := repository.CreateCountryRepository(db.DB)
	cityRepo := repository.CreateCityRepository(db.DB)

	// Создаем сервис словаря
	dictionaryService := services.CreateDictionaryService(
		translationRepo,
		countryRepo,
		cityRepo,
		redisService,
		cfg.Cache.TTL,
		loggerInstance,
	)

	// Создаем хендлер
	dictionaryHandler := handlers.CreateDictionaryHandler(dictionaryService, loggerInstance)

	// Создаем роутер
	r := gin.New()

	// Sentry: паники и 5xx уходят событиями; при пустом DSN остаётся
	// обычное восстановление паник с ответом 500.
	r.Use(sentrysdk.GinMiddleware("dictionary-service"))

	// Настраиваем роуты
	routes.SetupRoutes(r, dictionaryHandler, loggerInstance)

	// Создаем HTTP сервер с таймаутами
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Запускаем сервер в горутине
	go func() {
		loggerInstance.Infof("Dictionary service запущен на %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loggerInstance.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	loggerInstance.Info("Завершение работы сервера...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Закрываем подключения
	if redisService != nil {
		if err := redisService.Close(); err != nil {
			loggerInstance.Warnf("Ошибка закрытия Redis: %v", err)
		}
	}
	// Закрытие БД уже обрабатывается через defer выше

	if err := srv.Shutdown(ctx); err != nil {
		loggerInstance.Fatalf("Ошибка при завершении сервера: %v", err)
	}

	loggerInstance.Info("Сервер завершен")
}
