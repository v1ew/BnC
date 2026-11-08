package main

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
)

// Config описывает параметры партии «Быки и коровы»:
// длину загаданного кода и алфавит допустимых символов.
type Config struct {
	Length   int    // длина кода
	Alphabet []rune // уникальные символы, из которых собирается код
}

// presets — именованные алфавиты, доступные по имени через флаг -alphabet.
var presets = map[string]string{
	"digits":  "0123456789",
	"hex":     "0123456789abcdef",
	"letters": "abcdefghijklmnopqrstuvwxyz",
	"russian": "абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
}

// parseAlphabet превращает спецификацию алфавита в набор символов.
// Имя известного пресета (digits, hex, letters, russian) раскрывается
// в соответствующий набор; любая другая непустая строка считается явным
// перечислением символов алфавита.
func parseAlphabet(spec string) ([]rune, error) {
	if preset, ok := presets[strings.ToLower(spec)]; ok {
		return []rune(preset), nil
	}
	if spec == "" {
		return nil, fmt.Errorf("алфавит не задан")
	}
	return []rune(spec), nil
}

// validate проверяет, что по конфигурации можно играть: код имеет длину
// хотя бы один символ, алфавит не пуст и не короче кода, состоит из
// неповторяющихся символов и не содержит пробелов.
func (c Config) validate() error {
	if c.Length < 1 {
		return fmt.Errorf("длина кода должна быть не меньше 1")
	}
	if len(c.Alphabet) == 0 {
		return fmt.Errorf("алфавит пуст")
	}
	if len(c.Alphabet) < c.Length {
		return fmt.Errorf("алфавит слишком мал: в нём всего %d %s, а длина кода %d",
			len(c.Alphabet), symbolWord(len(c.Alphabet)), c.Length)
	}

	seen := make(map[rune]bool, len(c.Alphabet))
	for _, r := range c.Alphabet {
		if unicode.IsSpace(r) {
			return fmt.Errorf("пробельные символы не могут входить в алфавит: %q", r)
		}
		if seen[r] {
			return fmt.Errorf("символы в алфавите не должны повторяться: %q", r)
		}
		seen[r] = true
	}

	return nil
}

// newConfig собирает конфигурацию из длины кода и спецификации алфавита.
func newConfig(length int, alphabetSpec string) (Config, error) {
	alphabet, err := parseAlphabet(alphabetSpec)
	if err != nil {
		return Config{}, err
	}

	c := Config{Length: length, Alphabet: alphabet}
	if err := c.validate(); err != nil {
		return Config{}, err
	}

	return c, nil
}

// generateSecret возвращает случайный код длины c.Length, собранный из
// неповторяющихся символов алфавита. Конфигурация обязана быть валидной.
func (c Config) generateSecret(r *rand.Rand) string {
	runes := append([]rune(nil), c.Alphabet...)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes[:c.Length])
}

// validateGuess проверяет, что догадка — это c.Length неповторяющихся
// символов алфавита.
func (c Config) validateGuess(guess string) error {
	runes := []rune(guess)

	if len(runes) != c.Length {
		return fmt.Errorf("нужно ввести ровно %d %s", c.Length, symbolWord(c.Length))
	}

	inAlphabet := make(map[rune]bool, len(c.Alphabet))
	for _, r := range c.Alphabet {
		inAlphabet[r] = true
	}

	seen := make(map[rune]bool, c.Length)
	for _, r := range runes {
		if !inAlphabet[r] {
			return fmt.Errorf("символ %q не входит в алфавит игры", r)
		}
		if seen[r] {
			return fmt.Errorf("символы не должны повторяться: %q", r)
		}
		seen[r] = true
	}

	return nil
}

// countBullsCows считает быков (символ на своём месте) и коров (символ
// из кода на чужом месте). Код и догадка обязаны быть одинаковой длины;
// если это не так, лишние символы игнорируются.
func countBullsCows(secret, guess string) (bulls, cows int) {
	secretRunes := []rune(secret)
	guessRunes := []rune(guess)

	inSecret := make(map[rune]bool, len(secretRunes))
	for _, r := range secretRunes {
		inSecret[r] = true
	}

	n := len(secretRunes)
	if len(guessRunes) < n {
		n = len(guessRunes)
	}

	for i := 0; i < n; i++ {
		if guessRunes[i] == secretRunes[i] {
			bulls++
		} else if inSecret[guessRunes[i]] {
			cows++
		}
	}

	return bulls, cows
}

// symbolWord возвращает слово «символ» в форме, согласующейся с числом n:
// 1 символ, 2 символа, 5 символов.
func symbolWord(n int) string {
	switch {
	case n%10 == 1 && n%100 != 11:
		return "символ"
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 12 || n%100 > 14):
		return "символа"
	default:
		return "символов"
	}
}
