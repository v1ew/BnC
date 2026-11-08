package main

import (
	"math/rand"
	"strings"
	"testing"
)

// testConfigs — набор конфигураций для проверок свойств:
// разные длины и разные алфавиты, включая многобайтовые символы.
var testConfigs = []Config{
	{Length: 4, Alphabet: []rune("0123456789")},
	{Length: 6, Alphabet: []rune("0123456789abcdef")},
	{Length: 3, Alphabet: []rune("abcdef")},
	{Length: 2, Alphabet: []rune("#!@%&")},
	{Length: 4, Alphabet: []rune("абвгдеёжзийклмнопрстуфхцчшщъыьэюя")},
	{Length: 2, Alphabet: []rune("01")},
	{Length: 1, Alphabet: []rune("01")},
}

// --- Тесты конфигурации ---

func TestParseAlphabet(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		want    string
		wantErr bool
	}{
		{"пресет digits", "digits", "0123456789", false},
		{"имя пресета в любом регистре", "DIGITS", "0123456789", false},
		{"пресет hex", "hex", "0123456789abcdef", false},
		{"пресет letters", "letters", "abcdefghijklmnopqrstuvwxyz", false},
		{"пресет russian", "russian", "абвгдеёжзийклмнопрстуфхцчшщъыьэюя", false},
		{"явный набор символов", "abcd", "abcd", false},
		{"явный набор со знаками", "01#!", "01#!", false},
		{"пустая строка", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alphabet, err := parseAlphabet(tt.spec)

			if (err != nil) != tt.wantErr {
				t.Fatalf("parseAlphabet(%q) ошибка = %v, wantErr = %v",
					tt.spec, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if got := string(alphabet); got != tt.want {
				t.Errorf("parseAlphabet(%q) = %q, хотим %q", tt.spec, got, tt.want)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"классическая игра", Config{Length: 4, Alphabet: []rune("0123456789")}, false},
		{"алфавит ровно длины кода", Config{Length: 2, Alphabet: []rune("01")}, false},
		{"код из одного символа", Config{Length: 1, Alphabet: []rune("01")}, false},
		{"нулевая длина", Config{Length: 0, Alphabet: []rune("0123456789")}, true},
		{"отрицательная длина", Config{Length: -2, Alphabet: []rune("0123456789")}, true},
		{"пустой алфавит", Config{Length: 1, Alphabet: nil}, true},
		{"алфавит короче кода", Config{Length: 3, Alphabet: []rune("01")}, true},
		{"повтор символа в алфавите", Config{Length: 2, Alphabet: []rune("01210")}, true},
		{"пробел в алфавите", Config{Length: 2, Alphabet: []rune("01 2")}, true},
		{"табуляция в алфавите", Config{Length: 2, Alphabet: []rune("a\tbc")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() конфигурации %+v: ошибка = %v, wantErr = %v",
					tt.cfg, err, tt.wantErr)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	cfg, err := newConfig(6, "hex")
	if err != nil {
		t.Fatalf("newConfig(6, \"hex\") неожиданная ошибка: %v", err)
	}
	if cfg.Length != 6 {
		t.Errorf("длина кода = %d, хотим 6", cfg.Length)
	}
	if string(cfg.Alphabet) != "0123456789abcdef" {
		t.Errorf("алфавит = %q, хотим %q", string(cfg.Alphabet), "0123456789abcdef")
	}

	if _, err := newConfig(11, "digits"); err == nil {
		t.Errorf("newConfig(11, \"digits\") должен вернуть ошибку: в алфавите всего 10 цифр")
	}
	if _, err := newConfig(4, "01ab 3"); err == nil {
		t.Errorf("newConfig(4, \"01ab 3\") должен отклонить алфавит с пробелом")
	}
	if _, err := newConfig(3, "aba"); err == nil {
		t.Errorf("newConfig(3, \"aba\") должен отклонить алфавит с повтором символов")
	}
}

func TestSymbolWord(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "символ"},
		{21, "символ"},
		{101, "символ"},
		{2, "символа"},
		{4, "символа"},
		{22, "символа"},
		{5, "символов"},
		{10, "символов"},
		{11, "символов"},
		{12, "символов"},
		{111, "символов"},
	}

	for _, tt := range tests {
		if got := symbolWord(tt.n); got != tt.want {
			t.Errorf("symbolWord(%d) = %q, хотим %q", tt.n, got, tt.want)
		}
	}
}

// --- Тесты generateSecret ---

func TestGenerateSecret_Properties(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	for _, cfg := range testConfigs {
		allowed := make(map[rune]bool, len(cfg.Alphabet))
		for _, s := range cfg.Alphabet {
			allowed[s] = true
		}

		for i := 0; i < 1000; i++ {
			s := cfg.generateSecret(r)
			runes := []rune(s)

			if len(runes) != cfg.Length {
				t.Fatalf("длина кода = %d, ожидалось %d (код: %q, конфигурация: %+v)",
					len(runes), cfg.Length, s, cfg)
			}

			seen := make(map[rune]bool, cfg.Length)
			for _, ch := range runes {
				if !allowed[ch] {
					t.Fatalf("символ %q вне алфавита (код: %q, конфигурация: %+v)",
						ch, s, cfg)
				}
				if seen[ch] {
					t.Fatalf("символы кода повторяются (код: %q, конфигурация: %+v)",
						s, cfg)
				}
				seen[ch] = true
			}
		}
	}
}

func TestGenerateSecret_Variety(t *testing.T) {
	// За 10 генераций должно встретиться хотя бы два разных кода.
	cfg := Config{Length: 4, Alphabet: []rune("0123456789")}
	r := rand.New(rand.NewSource(42))

	seen := make(map[string]bool)
	for i := 0; i < 10; i++ {
		seen[cfg.generateSecret(r)] = true
	}

	if len(seen) < 2 {
		t.Fatalf("за 10 генераций получен один и тот же код: %v", seen)
	}
}

func TestGenerateSecret_TinyAlphabetCoversAllCodes(t *testing.T) {
	// Когда алфавит ровно длины кода, возможны только перестановки —
	// за 100 генераций обязаны встретиться обе.
	cfg := Config{Length: 2, Alphabet: []rune("01")}
	r := rand.New(rand.NewSource(42))

	seen := map[string]bool{"01": false, "10": false}
	for i := 0; i < 100; i++ {
		seen[cfg.generateSecret(r)] = true
	}

	for code, met := range seen {
		if !met {
			t.Errorf("код %q ни разу не выпал за 100 генераций", code)
		}
	}
}

func TestGenerateSecret_AllowsLeadingZero(t *testing.T) {
	// Код — это последовательность символов, а не число: первая позиция
	// может быть любой, включая ноль.
	cfg := Config{Length: 4, Alphabet: []rune("0123456789")}
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 1000; i++ {
		if strings.HasPrefix(cfg.generateSecret(r), "0") {
			return
		}
	}

	t.Fatalf("за 1000 генераций ни один код не начался с нуля")
}

func TestGenerateSecret_IsValidGuess(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	for _, cfg := range testConfigs {
		for i := 0; i < 100; i++ {
			s := cfg.generateSecret(r)
			if err := cfg.validateGuess(s); err != nil {
				t.Fatalf("сгенерированный код %q не проходит валидацию: %v (конфигурация: %+v)",
					s, err, cfg)
			}
		}
	}
}

// --- Тесты validateGuess ---

func TestValidateGuess(t *testing.T) {
	digits := Config{Length: 4, Alphabet: []rune("0123456789")}
	letters := Config{Length: 4, Alphabet: []rune("abcdefgh")}
	russian := Config{Length: 4, Alphabet: []rune("абвгдеёж")}

	tests := []struct {
		name    string
		cfg     Config
		input   string
		wantErr bool
	}{
		{"цифры: валидный код", digits, "1234", false},
		{"цифры: ноль в начале", digits, "0123", false},
		{"цифры: ноль в конце", digits, "9870", false},
		{"цифры: пустая строка", digits, "", true},
		{"цифры: короче", digits, "123", true},
		{"цифры: длиннее", digits, "12345", true},
		{"цифры: буква внутри", digits, "12a4", true},
		{"цифры: буква в начале", digits, "a123", true},
		{"цифры: пробел внутри", digits, "12 4", true},
		{"цифры: только пробелы", digits, "    ", true},
		{"цифры: повторяющиеся цифры", digits, "1123", true},
		{"цифры: все одинаковые", digits, "1111", true},
		{"цифры: минус", digits, "-123", true},
		{"цифры: спецсимволы", digits, "12#4", true},
		{"буквы: валидный код", letters, "abcd", false},
		{"буквы: цифра вне алфавита", letters, "ab1d", true},
		{"буквы: верхний регистр вне алфавита", letters, "Abcd", true},
		{"буквы: повтор", letters, "abca", true},
		{"русские буквы: валидный код", russian, "беда", false},
		{"русские буквы: латиница вне алфавита", russian, "abcd", true},
		{"русские буквы: длиннее по символам", russian, "абвгд", true},
		{"русские буквы: короче по символам", russian, "абв", true},
		{"русские буквы: повтор", russian, "абва", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validateGuess(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGuess(%q) ошибка = %v, wantErr = %v",
					tt.input, err, tt.wantErr)
			}
		})
	}
}

// --- Тесты countBullsCows ---

func TestCountBullsCows(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		guess     string
		wantBulls int
		wantCows  int
	}{
		{"точное совпадение", "1234", "1234", 4, 0},
		{"все коровы (обратный порядок)", "1234", "4321", 0, 4},
		{"нет совпадений", "1234", "5678", 0, 0},
		{"два быка и две коровы", "1234", "1243", 2, 2},
		{"три быка", "1234", "1235", 3, 0},
		{"один бык и две коровы", "1234", "1523", 1, 2},
		{"три быка кроме последней", "1234", "1294", 3, 0},
		{"ноль в коде, все коровы", "0123", "3210", 0, 4},
		{"ноль в коде, быки и коровы", "0123", "0321", 2, 2},
		{"три коровы, ноль быков", "5678", "7564", 0, 3},
		{"один бык и две коровы (2)", "5678", "7654", 1, 2},
		{"латиница: точное совпадение", "abcd", "abcd", 4, 0},
		{"латиница: все коровы", "abcd", "dcba", 0, 4},
		{"кириллица: все коровы", "абвг", "гвба", 0, 4},
		{"кириллица: два быка и две коровы", "абвг", "агвб", 2, 2},
		{"код из одного символа: бык", "7", "7", 1, 0},
		{"код из одного символа: мимо", "7", "3", 0, 0},
		{"разная длина: хвост догадки игнорируется", "ab", "abc", 2, 0},
		{"разная длина: короткая догадка", "abc", "ab", 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bulls, cows := countBullsCows(tt.secret, tt.guess)
			if bulls != tt.wantBulls || cows != tt.wantCows {
				t.Errorf("countBullsCows(%q, %q) = (%d, %d), хотим (%d, %d)",
					tt.secret, tt.guess, bulls, cows, tt.wantBulls, tt.wantCows)
			}
		})
	}
}

func TestCountBullsCows_Identity(t *testing.T) {
	// Для любого сгенерированного кода сравнение с самим собой
	// даёт ровно Length быков и ноль коров.
	r := rand.New(rand.NewSource(42))

	for _, cfg := range testConfigs {
		for i := 0; i < 100; i++ {
			s := cfg.generateSecret(r)
			bulls, cows := countBullsCows(s, s)

			if bulls != cfg.Length || cows != 0 {
				t.Fatalf("countBullsCows(%q, %q) = (%d, %d), хотим (%d, 0)",
					s, s, bulls, cows, cfg.Length)
			}
		}
	}
}

func TestCountBullsCows_ReversedSecretSumsToLength(t *testing.T) {
	// Перестановка символов кода: быки + коровы всегда равны длине кода.
	r := rand.New(rand.NewSource(42))

	for _, cfg := range testConfigs {
		for i := 0; i < 100; i++ {
			secret := cfg.generateSecret(r)
			reversed := reverseRunes(secret)

			bulls, cows := countBullsCows(secret, reversed)
			if bulls+cows != cfg.Length {
				t.Fatalf("countBullsCows(%q, %q) = (%d, %d), сумма должна быть %d",
					secret, reversed, bulls, cows, cfg.Length)
			}
		}
	}
}

// reverseRunes переворачивает строку по символам (не по байтам).
func reverseRunes(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
