package handlers

// SDL-315: инварианты sanitiseQueryString.
//
// Всё, что проходит через этот хелпер, уезжает в запрос к Postgres, а
// Postgres не принимает NULL-байты и невалидный UTF-8 в TEXT — такой ввод
// возвращается пользователю как 500. Поэтому единственный контракт хелпера:
// на выходе строка, которую Postgres примет, при любом вводе.

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSanitiseQueryStringRemovesNullAndControlChars(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"NULL-байт вырезается", "A\x00B", "AB"},
		{"управляющие символы вырезаются", "A\x01\x02\x1fB", "AB"},
		{"таб сохраняется", "A\tB", "A\tB"},
		{"обычный текст не трогается", "KZ", "KZ"},
		{"кириллица не трогается", "Алматы", "Алматы"},
		{"пустая строка", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := sanitiseQueryString(tc.in); got != tc.want {
				t.Errorf("sanitiseQueryString(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestSanitiseQueryStringAlwaysReturnsValidUTF8 — главный инвариант.
// Невалидный UTF-8 на входе (одиночный continuation-байт, оборванная
// последовательность) обязан стать валидным на выходе, иначе Postgres
// отвечает ошибкой и пользователь получает 500 на кривой query-параметр.
func TestSanitiseQueryStringAlwaysReturnsValidUTF8(t *testing.T) {
	inputs := []string{
		"\xc2",                 // оборванная 2-байтовая последовательность
		"\xff\xfe",             // байты, невозможные в UTF-8
		"KZ\xc3",               // валидный префикс + обрывок
		"\xf3\x84\x82\xa3\x1d", // из фаззера, ронявшего /cities
		strings.Repeat("я", 300),
		strings.Repeat("a", 300),
		// Ровно та длина, на которой байтовый срез разрубал бы руну:
		// 199 однобайтовых + двухбайтовая на границе 200.
		strings.Repeat("a", 199) + "я",
	}
	for _, in := range inputs {
		got := sanitiseQueryString(in)
		if !utf8.ValidString(got) {
			t.Errorf("sanitiseQueryString(%q) вернул невалидный UTF-8: %q", in, got)
		}
		if strings.ContainsRune(got, 0) {
			t.Errorf("sanitiseQueryString(%q) оставил NULL-байт: %q", in, got)
		}
	}
}

// TestSanitiseQueryStringCapsLength — потолок длины остаётся в силе, и
// обрезка не создаёт невалидный UTF-8 (регресс на байтовый срез).
func TestSanitiseQueryStringCapsLength(t *testing.T) {
	long := strings.Repeat("я", 300) // 600 байт
	got := sanitiseQueryString(long)

	if len(got) > maxQueryStringLen {
		t.Errorf("длина %d превышает потолок %d", len(got), maxQueryStringLen)
	}
	if !utf8.ValidString(got) {
		t.Errorf("обрезка разрубила руну, получился невалидный UTF-8: %q", got)
	}
	// Обрезка по границе руны может снять до 1 байта сверх потолка;
	// больше — значит режем лишнее.
	if len(got) < maxQueryStringLen-1 {
		t.Errorf("обрезано слишком много: %d байт, ожидалось около %d", len(got), maxQueryStringLen)
	}
}
