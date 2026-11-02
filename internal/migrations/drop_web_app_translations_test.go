package migrations

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SDL-313. Тест исполняется против НАСТОЯЩЕГО Postgres, а не мока: проверяемое
// поведение — идемпотентность DELETE и то, что чужие строки уцелели, — целиком
// про поведение БД, и на моке доказывается только текст запроса.
//
// DSN берётся из SDL313_TEST_DSN. Без переменной тест пропускается: в пайплайне
// dictionary-service нет service-контейнера с Postgres (это записанный долг,
// см. закрывающий коммент задачи). При приёмке 21 Aug 2026 тест прогонялся
// вручную против одноразовой БД на тест-кластере CNPG.
func TestDropWebAppTranslationsIsIdempotent(t *testing.T) {
	dsn := os.Getenv("SDL313_TEST_DSN")
	if dsn == "" {
		t.Skip("SDL313_TEST_DSN не задан — нужен живой Postgres")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("подключение к Postgres: %v", err)
	}

	if err := db.Exec(`DROP TABLE IF EXISTS translations`).Error; err != nil {
		t.Fatalf("подготовка: %v", err)
	}
	if err := db.Exec(`CREATE TABLE translations (id serial primary key, app text, key text, value text)`).Error; err != nil {
		t.Fatalf("создание фикстуры: %v", err)
	}
	if err := db.Exec(`INSERT INTO translations (app, key, value) VALUES
		('web-app','a','1'), ('web-app','b','2'), ('web-app','c','3'),
		('telegram-bot','d','4'), (NULL,'e','5')`).Error; err != nil {
		t.Fatalf("наполнение фикстуры: %v", err)
	}

	removed, err := DropWebAppTranslations(db)
	if err != nil {
		t.Fatalf("первый прогон: %v", err)
	}
	if removed != 3 {
		t.Errorf("первый прогон удалил %d строк, ожидалось 3", removed)
	}

	// Второй прогон — тот самый рестарт пода: он обязан быть no-op, а не
	// ошибкой и не повторным удалением.
	removed, err = DropWebAppTranslations(db)
	if err != nil {
		t.Fatalf("повторный прогон: %v", err)
	}
	if removed != 0 {
		t.Errorf("повторный прогон удалил %d строк, ожидался 0 — миграция не идемпотентна", removed)
	}

	var webApp, others int64
	db.Raw(`SELECT count(*) FROM translations WHERE app = 'web-app'`).Scan(&webApp)
	db.Raw(`SELECT count(*) FROM translations WHERE app IS DISTINCT FROM 'web-app'`).Scan(&others)
	if webApp != 0 {
		t.Errorf("после чистки осталось %d строк app='web-app'", webApp)
	}
	// Соседи должны уцелеть: и другое приложение, и строки с app IS NULL —
	// на проде их 240 и 78 соответственно, и они к SDL-311 отношения не имеют.
	if others != 2 {
		t.Errorf("чистка задела чужие строки: осталось %d из 2", others)
	}

	db.Exec(`DROP TABLE translations`)
}
