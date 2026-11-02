package migrations

import (
	"strings"
	"testing"
)

// Сид-корпус читается через go:embed, поэтому опечатка в пути или пустой
// файл не ломают сборку — RunSeeds лишь пишет Warn и идёт дальше (сидер
// non-fatal by design). Тест закрывает эту слепую зону: он исполняет тот же
// обход, что и RunSeeds, и падает там, где рантайм промолчал бы.
func TestEmbeddedSeedsAreReadableAndNonEmpty(t *testing.T) {
	entries, err := seedFiles.ReadDir("seeds")
	if err != nil {
		t.Fatalf("seeds недоступны через embed: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("корпус сидов пуст — RunSeeds стал бы no-op, не сказав ни слова")
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			t.Errorf("в seeds/ лежит не-.sql запись %q — RunSeeds её молча пропустит", e.Name())
			continue
		}
		data, err := seedFiles.ReadFile("seeds/" + e.Name())
		if err != nil {
			t.Errorf("seeds/%s не читается: %v", e.Name(), err)
			continue
		}
		if strings.TrimSpace(string(data)) == "" {
			t.Errorf("seeds/%s пуст — RunSeeds его пропустит, а справочник останется без строк", e.Name())
		}
	}
}

// SDL-311: строки app='web-app' сидились для фронта, у которого нет и не было
// читателя (каталог next-intl живёт в markerloop-app/apps/admin/src/messages
// и в БД не ходит). Сид удалён; тест не даёт ему вернуться незамеченным —
// иначе расхождение всплывёт только через ~900 мёртвых строк в проде.
func TestNoWebAppRowsAreSeeded(t *testing.T) {
	entries, err := seedFiles.ReadDir("seeds")
	if err != nil {
		t.Fatalf("seeds недоступны через embed: %v", err)
	}

	for _, e := range entries {
		data, err := seedFiles.ReadFile("seeds/" + e.Name())
		if err != nil {
			continue
		}
		if strings.Contains(string(data), "'web-app'") {
			t.Errorf("seeds/%s снова сидит app='web-app' — у этих строк нет читателя (SDL-311)", e.Name())
		}
	}
}
