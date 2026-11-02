package dictionary

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"dictionary-service/internal/models"
)

// SDL-429. Гард против класса дефекта, а не против одного его случая.
//
// GetError на промах по коду возвращает не ошибку, а запись INTERNAL_SERVER_ERROR,
// поэтому отсутствие строки в словаре ничего не роняет: клиент молча получает
// «внутренняя ошибка сервера» вместо своей ошибки, а вызывающая сторона (gw-service
// в sendForbiddenError) видит err=nil и считает ответ валидным. Именно так код
// FORBIDDEN — объявленный, но не заведённый — полгода отдавал 403 с телом про
// внутреннюю ошибку. Вместе с ним не заведены были CONFLICT и VALIDATION_ERROR.
//
// Список кодов берётся разбором исходника models, а не копией в тесте: копия
// разъедется с моделью ровно так же, как разъехался словарь, и гард снова начнёт
// проверять сам себя. Добавили новую константу ErrorCode — тест обязан упасть,
// пока для неё не заведены строки на всех трёх языках.
func TestEveryDeclaredErrorCodeHasDictionaryEntry(t *testing.T) {
	declared := parseDeclaredErrorCodes(t, "../../internal/models/dictionary.go")
	if len(declared) == 0 {
		t.Fatal("в models не найдено ни одной константы ErrorCode — сломан разбор, а не словарь")
	}

	languages := []models.LanguageCode{models.LanguageRU, models.LanguageKK, models.LanguageEN}

	for _, name := range declared {
		code := models.ErrorCode(name)
		entry, ok := ErrorDictionary[code]
		if !ok {
			t.Errorf("код %s объявлен в models, но записи в ErrorDictionary нет: GetError молча отдаст INTERNAL_SERVER_ERROR", name)
			continue
		}
		for _, lang := range languages {
			localized, ok := entry[lang]
			if !ok {
				t.Errorf("код %s: нет строки для языка %s", name, lang)
				continue
			}
			if localized.Message == "" {
				t.Errorf("код %s, язык %s: пустое Message", name, lang)
			}
			if localized.Code != code {
				t.Errorf("код %s, язык %s: в записи проставлен чужой Code=%s", name, lang, localized.Code)
			}
		}
	}
}

// parseDeclaredErrorCodes достаёт значения всех констант типа ErrorCode из исходника
// моделей. Нужны именно значения ("FORBIDDEN"), а не имена констант (ErrorForbidden):
// ключом словаря и содержимым ответа служит значение.
func parseDeclaredErrorCodes(t *testing.T, path string) []string {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatalf("разбор %s: %v", path, err)
	}

	var codes []string
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		// Тип пишется у первой спецификации блока и наследуется остальными.
		currentType := ""
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if ident, ok := valueSpec.Type.(*ast.Ident); ok {
				currentType = ident.Name
			} else if valueSpec.Type != nil {
				currentType = ""
			}
			if currentType != "ErrorCode" {
				continue
			}
			for _, value := range valueSpec.Values {
				lit, ok := value.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				codes = append(codes, lit.Value[1:len(lit.Value)-1])
			}
		}
	}
	return codes
}
