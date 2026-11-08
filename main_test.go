package main

import (
	"math/rand"
	"strings"
	"testing"
)

// --- Тесты для generateSecret ---

func TestGenerateSecret_Properties(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 1000; i++ {
		s := generateSecret(r)

		if len(s) != secretLength {
			t.Fatalf("длина числа = %d, ожидалось %d (число: %q)", len(s), secretLength, s)
		}

		if s[0] == '0' {
			t.Fatalf("первая цифра не должна быть нулём (число: %q)", s)
		}

		seen := make(map[rune]bool, secretLength)
		for _, ch := range s {
			if ch < '0' || ch > '9' {
				t.Fatalf("в числе есть нецифровой символ: %q", s)
			}
			if seen[ch] {
				t.Fatalf("в числе есть повторяющиеся цифры: %q", s)
			}
			seen[ch] = true
		}
	}
}

func TestGenerateSecret_Variety(t *testing.T) {
	// За 20 генераций должно встретиться хотя бы 2 разных числа.
	r := rand.New(rand.NewSource(42))
	seen := make(map[string]bool, 10)

	for i := 0; i < 10; i++ {
		seen[generateSecret(r)] = true
	}

	if len(seen) < 2 {
		t.Fatalf("за 10 генераций получено одно и то же число")
	}
}

func TestGenerateSecret_IsValidGuess(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		s := generateSecret(r)
		if err := validateGuess(s); err != nil {
			t.Fatalf("сгенерированное число %q не проходит валидацию: %v", s, err)
		}
	}
}

// --- Тесты для validateGuess ---

func TestValidateGuess(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"валидное число", "1234", false},
		{"с нулём в начале", "0123", false},
		{"с нулём в конце", "9870", false},
		{"пустая строка", "", true},
		{"короче 4 цифр", "123", true},
		{"длиннее 4 цифр", "12345", true},
		{"буква внутри", "12a4", true},
		{"буква в начале", "a123", true},
		{"пробел внутри", "12 4", true},
		{"только пробелы", "    ", true},
		{"повторяющиеся цифры", "1123", true},
		{"все одинаковые", "1111", true},
		{"отрицательное", "-123", true},
		{"спецсимволы", "12#4", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGuess(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateGuess(%q) ошибка = %v, wantErr = %v",
					tt.input, err, tt.wantErr)
			}
		})
	}
}

// --- Тесты для countBullsCows ---

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
		{"ноль в числе, все коровы", "0123", "3210", 0, 4},
		{"ноль в числе, быки и коровы", "0123", "0321", 2, 2},
		{"три коровы, ноль быков", "5678", "7564", 0, 3},
		{"один бык и две коровы", "5678", "7654", 1, 2},
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
	// Для любого сгенерированного числа быков должно быть ровно secretLength, коров — 0.
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		s := generateSecret(r)
		bulls, cows := countBullsCows(s, s)

		if bulls != secretLength || cows != 0 {
			t.Fatalf("countBullsCows(%q, %q) = (%d, %d), хотим (%d, 0)",
				s, s, bulls, cows, secretLength)
		}
	}
}

func TestCountBullsCows_AllDigitsFromSecret(t *testing.T) {
	// Если переставить цифры секрета, быки+коровы всегда = secretLength.
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		secret := generateSecret(r)

		// Переворачиваем строку — это всегда перестановка тех же цифр.
		reversed := reverseString(secret)

		bulls, cows := countBullsCows(secret, reversed)
		if bulls+cows != secretLength {
			t.Fatalf("countBullsCows(%q, %q) = (%d, %d), сумма должна быть %d",
				secret, reversed, bulls, cows, secretLength)
		}
	}
}

func reverseString(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}

