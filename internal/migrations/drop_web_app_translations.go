package migrations

import (
	"gorm.io/gorm"
)

// DropWebAppTranslations удаляет строки translations с app='web-app'.
//
// SDL-311 убрала сид этих строк: у них нет и не было читателя — каталог
// админки живёт в markerloop-app/apps/admin/src/messages и в БД не ходит.
// Но сидер работает через ON CONFLICT и никогда не делает DELETE, поэтому
// уже засеянные строки удаление сида пережили: на 21 Aug 2026 их ровно 900
// и на тесте, и на проде. Эта функция их и убирает (чистка утверждена
// владельцем 21 Aug 2026, SDL-313).
//
// Идемпотентна: повторный запуск удаляет 0 строк и возвращает nil. Поэтому
// её безопасно звать на каждом старте пода — а звать на каждом старте нужно,
// иначе инвариант «после старта строк app='web-app' нет» держался бы только
// до следующего восстановления базы из старого бэкапа.
//
// Откат — git-история сида: строки восстанавливаются повторным сидом из
// коммита до SDL-311.
func DropWebAppTranslations(db *gorm.DB) (int64, error) {
	res := db.Exec(`DELETE FROM translations WHERE app = 'web-app'`)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
