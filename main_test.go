package main

import (
	"math/rand"
	"strings"
	"testing"
)

// knownSecret возвращает код, который runGame сгенерирует для cfg при сиде 42:
// генерация детерминирована, одинаковый сид даёт одинаковый код.
func knownSecret(cfg Config) string {
	return cfg.generateSecret(rand.New(rand.NewSource(42)))
}

// playRound прогоняет партию с сидом 42 и возвращает весь вывод игры.
func playRound(t *testing.T, cfg Config, input string) string {
	t.Helper()

	var out strings.Builder
	runGame(cfg, rand.New(rand.NewSource(42)), strings.NewReader(input), &out)
	return out.String()
}

func TestRunGame_Win(t *testing.T) {
	cfg := Config{Length: 4, Alphabet: []rune("0123456789")}
	secret := knownSecret(cfg)

	output := playRound(t, cfg, "abcd\n1234\n"+secret+"\n")

	if !strings.Contains(output, "Ошибка:") {
		t.Errorf("вывод должен содержать сообщение об ошибке:\n%s", output)
	}
	if !strings.Contains(output, "символ 'a' не входит в алфавит игры") {
		t.Errorf("вывод должен объяснить, что символ вне алфавита:\n%s", output)
	}
	if !strings.Contains(output, "Поздравляю! Вы угадали код: "+secret) {
		t.Errorf("вывод должен содержать сообщение о победе:\n%s", output)
	}
	if !strings.Contains(output, "Ваш вариант: ") {
		t.Errorf("вывод должен содержать приглашение ко вводу:\n%s", output)
	}
}

func TestRunGame_HintsAndExit(t *testing.T) {
	cfg := Config{Length: 4, Alphabet: []rune("0123456789")}
	secret := knownSecret(cfg)

	// Меняем местами первые два символа кода: такая догадка всегда
	// даёт ровно двух быков и двух коров.
	runes := []rune(secret)
	runes[0], runes[1] = runes[1], runes[0]
	guess := string(runes)

	output := playRound(t, cfg, guess+"\nexit\n")

	if !strings.Contains(output, "Быков: 2\nКоров: 2") {
		t.Errorf("вывод должен содержать подсказку «Быков: 2» и «Коров: 2»:\n%s", output)
	}
	if !strings.Contains(output, "Игра завершена.") {
		t.Errorf("команда exit должна завершать игру:\n%s", output)
	}
}

func TestRunGame_ValidationErrors(t *testing.T) {
	cfg := Config{Length: 4, Alphabet: []rune("0123456789")}

	output := playRound(t, cfg, "1123\nab\nexit\n")

	if !strings.Contains(output, "символы не должны повторяться") {
		t.Errorf("повтор символов должен вызывать ошибку:\n%s", output)
	}
	if !strings.Contains(output, "нужно ввести ровно 4 символа") {
		t.Errorf("короткая догадка должна вызывать ошибку:\n%s", output)
	}
}

func TestRunGame_ExitCommands(t *testing.T) {
	cfg := Config{Length: 4, Alphabet: []rune("0123456789")}

	for _, cmd := range []string{"exit", "quit"} {
		output := playRound(t, cfg, cmd+"\n")
		if !strings.Contains(output, "Игра завершена.") {
			t.Errorf("команда %q должна завершать игру:\n%s", cmd, output)
		}
	}
}

func TestRunGame_EOF(t *testing.T) {
	cfg := Config{Length: 4, Alphabet: []rune("0123456789")}

	output := playRound(t, cfg, "")

	if !strings.Contains(output, "Ввод завершён.") {
		t.Errorf("конец ввода должен завершать игру:\n%s", output)
	}
}

func TestRunGame_HexConfig(t *testing.T) {
	cfg := Config{Length: 5, Alphabet: []rune("0123456789abcdef")}
	secret := knownSecret(cfg)

	output := playRound(t, cfg, "abc\n"+secret+"\n")

	if !strings.Contains(output, "нужно ввести ровно 5 символов") {
		t.Errorf("короткая догадка должна вызывать ошибку:\n%s", output)
	}
	if !strings.Contains(output, "Алфавит: 0123456789abcdef") {
		t.Errorf("стартовое сообщение должно показывать алфавит:\n%s", output)
	}
	if !strings.Contains(output, "Поздравляю! Вы угадали код: "+secret) {
		t.Errorf("вывод должен содержать сообщение о победе:\n%s", output)
	}
}

func TestRunGame_RussianConfig(t *testing.T) {
	cfg := Config{Length: 3, Alphabet: []rune("абвгдеёж")}
	secret := knownSecret(cfg)

	output := playRound(t, cfg, secret+"\n")

	if !strings.Contains(output, "Поздравляю! Вы угадали код: "+secret) {
		t.Errorf("многобайтовый алфавит должен работать от начала до конца:\n%s", output)
	}
}
