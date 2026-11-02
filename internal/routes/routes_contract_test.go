package routes_test

// SDL-315: контракт группы /api/v1/dictionary как её видит шлюз.
//
// Шлюз (gw-service) публикует эти роуты наружу и проксирует путь запроса
// в сервис как есть. Если метод роута здесь разойдётся с методом на
// шлюзе, эндпоинт станет недостижим МОЛЧА: правильный метод получит
// Gin-дефолт 404 на шлюзе, метод шлюза — Gin-дефолт 404 от этого сервиса.
// Ровно так /v1/dictionary/cities пролежал недоступным на обоих контурах
// (шлюз держал POST, сервис — GET).
//
// У шлюза лежит зеркало этой же таблицы:
// gw-service/internal/routes/dictionary_contract_test.go
// (переменная dictionaryUpstreamRoutes). Этот тест сверяет таблицу с
// ФАКТИЧЕСКИМ роутером сервиса, поэтому изменение роутов здесь роняет
// сборку здесь же — до того, как расхождение доедет до контура. При
// падении правь ОБЕ копии: эту и зеркало на шлюзе.

import (
	"io"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"dictionary-service/internal/routes"
)

const dictionaryPrefix = "/api/v1/dictionary"

// dictionaryContract — ожидаемый набор роутов группы /api/v1/dictionary.
// Держать синхронно с dictionaryUpstreamRoutes в gw-service.
var dictionaryContract = map[string]string{
	"/error":        http.MethodPost,
	"/errors":       http.MethodGet,
	"/translations": http.MethodPost,
	// Алиас для go-sdk (сервисы ходят напрямую); на шлюз не публикуется.
	"/words":     http.MethodPost,
	"/countries": http.MethodGet,
	"/cities":    http.MethodGet,
	"/languages": http.MethodGet,
}

// actualDictionaryRoutes — реальные роуты сервиса под /api/v1/dictionary.
// Хендлер можно не конструировать: SetupRoutes только берёт значения
// методов у указателя и не разыменовывает его при регистрации.
func actualDictionaryRoutes(t *testing.T) map[string][]string {
	t.Helper()
	gin.SetMode(gin.TestMode)

	log := logrus.New()
	log.SetOutput(io.Discard)

	r := gin.New()
	routes.SetupRoutes(r, nil, log)

	got := make(map[string][]string)
	for _, ri := range r.Routes() {
		if !strings.HasPrefix(ri.Path, dictionaryPrefix+"/") {
			continue
		}
		path := strings.TrimPrefix(ri.Path, dictionaryPrefix)
		got[path] = append(got[path], ri.Method)
	}
	for p := range got {
		sort.Strings(got[p])
	}
	return got
}

// TestDictionaryRouteTableMatchesContract — набор роутов сервиса обязан
// совпадать с контрактом по методу и пути, в обе стороны.
func TestDictionaryRouteTableMatchesContract(t *testing.T) {
	got := actualDictionaryRoutes(t)

	for path, wantMethod := range dictionaryContract {
		methods, ok := got[path]
		if !ok {
			t.Errorf("контракт ожидает %s %s%s, но сервис такой роут не регистрирует — шлюз, публикующий его, будет получать 404 (удали роут из контракта и из зеркала в gw-service, если он снят намеренно)",
				wantMethod, dictionaryPrefix, path)
			continue
		}
		if len(methods) != 1 || methods[0] != wantMethod {
			t.Errorf("метод %s%s разошёлся с контрактом: сервис отдаёт %s, контракт ожидает %s — обнови контракт ЗДЕСЬ и зеркало dictionaryUpstreamRoutes в gw-service/internal/routes/dictionary_contract_test.go, иначе шлюз останется на старом методе и эндпоинт станет недостижим",
				dictionaryPrefix, path, strings.Join(methods, ","), wantMethod)
		}
	}

	for path, methods := range got {
		if _, ok := dictionaryContract[path]; !ok {
			t.Errorf("сервис регистрирует %s %s%s, которого нет в контракте — добавь его сюда и реши, публиковать ли его на шлюзе (gw-service/internal/routes/dictionary_contract_test.go)",
				strings.Join(methods, ","), dictionaryPrefix, path)
		}
	}
}
