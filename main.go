package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const secretLength = 4

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	secret := generateSecret(r)

	fmt.Println("Игра «Быки и коровы»")
	fmt.Println("Я загадал четырёхзначное число без повторяющихся цифр.")
	fmt.Println("Введите 4 цифры или 'exit' для выхода.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Ваш вариант: ")

		if !scanner.Scan() {
			fmt.Println("\nВвод завершён.")
			return
		}

		guess := strings.TrimSpace(scanner.Text())

		if guess == "exit" || guess == "quit" {
			fmt.Println("Игра завершена.")
			return
		}

		if err := validateGuess(guess); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			continue
		}

		if guess == secret {
			fmt.Printf("Поздравляю! Вы угадали число: %s\n", secret)
			return
		}

		bulls, cows := countBullsCows(secret, guess)

		if bulls == 0 && cows == 0 {
			fmt.Println("Совпадений нет.")
		} else {
			if bulls > 0 {
				fmt.Printf("Быков: %d\n", bulls)
			}
			if cows > 0 {
				fmt.Printf("Коров: %d\n", cows)
			}
		}
	}
}

// generateSecret генерирует четырёхзначное число без повторяющихся цифр.
// Первая цифра не может быть нулём.
func generateSecret(r *rand.Rand) string {
	used := make([]bool, 10)
	result := make([]byte, 0, secretLength)

	first := r.Intn(9) + 1 // 1..9
	used[first] = true
	result = append(result, byte('0'+first))

	for len(result) < secretLength {
		d := r.Intn(10)

		if used[d] {
			continue
		}

		used[d] = true
		result = append(result, byte('0'+d))
	}

	return string(result)
}

// validateGuess проверяет, что введено ровно 4 цифры и они не повторяются.
func validateGuess(s string) error {
	if len(s) != secretLength {
		return fmt.Errorf("нужно ввести ровно 4 цифры")
	}

	seen := make(map[byte]bool, secretLength)

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch < '0' || ch > '9' {
			return fmt.Errorf("можно вводить только цифры")
		}

		if seen[ch] {
			return fmt.Errorf("цифры не должны повторяться")
		}

		seen[ch] = true
	}

	return nil
}

// countBullsCows считает быков и коров.
func countBullsCows(secret, guess string) (bulls int, cows int) {
	secretDigits := make(map[byte]bool, secretLength)

	for i := 0; i < len(secret); i++ {
		secretDigits[secret[i]] = true
	}

	for i := 0; i < secretLength; i++ {
		if guess[i] == secret[i] {
			bulls++
		} else if secretDigits[guess[i]] {
			cows++
		}
	}

	return bulls, cows
}

