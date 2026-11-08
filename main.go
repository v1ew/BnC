// Командная игра «Быки и коровы».
//
// Программа загадывает код из неповторяющихся символов и на каждую догадку
// отвечает, сколько символов стоит на своих местах (быки) и сколько входит
// в код, но на чужих местах (коровы).
//
// Параметры запуска:
//
//	-length N   длина загаданного кода, N >= 1 (по умолчанию 4)
//	-alphabet S алфавит кода: имя пресета (digits, hex, letters, russian)
//	            или явный набор неповторяющихся символов (по умолчанию digits)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	length := flag.Int("length", 4, "длина загаданного кода (не меньше 1)")
	alphabet := flag.String("alphabet", "digits",
		"алфавит кода: имя пресета (digits, hex, letters, russian) или явный набор символов")
	flag.Parse()

	cfg, err := newConfig(*length, *alphabet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	runGame(cfg, r, os.Stdin, os.Stdout)
}

// runGame проводит одну партию: загадывает код по конфигурации cfg и читает
// догадки из in, пока игрок не угадает, не введёт 'exit'/'quit' или не
// закончится ввод.
func runGame(cfg Config, r *rand.Rand, in io.Reader, out io.Writer) {
	secret := cfg.generateSecret(r)

	fmt.Fprintln(out, "Игра «Быки и коровы»")
	fmt.Fprintf(out, "Я загадал код длиной %d без повторов. Алфавит: %s\n",
		cfg.Length, string(cfg.Alphabet))
	fmt.Fprintf(out, "Введите код длиной %d или 'exit' для выхода.\n", cfg.Length)

	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, "Ваш вариант: ")

		if !scanner.Scan() {
			fmt.Fprintln(out, "\nВвод завершён.")
			return
		}

		guess := strings.TrimSpace(scanner.Text())

		if guess == "exit" || guess == "quit" {
			fmt.Fprintln(out, "Игра завершена.")
			return
		}

		if err := cfg.validateGuess(guess); err != nil {
			fmt.Fprintf(out, "Ошибка: %v\n", err)
			continue
		}

		if guess == secret {
			fmt.Fprintf(out, "Поздравляю! Вы угадали код: %s\n", secret)
			return
		}

		bulls, cows := countBullsCows(secret, guess)

		if bulls == 0 && cows == 0 {
			fmt.Fprintln(out, "Совпадений нет.")
		} else {
			if bulls > 0 {
				fmt.Fprintf(out, "Быков: %d\n", bulls)
			}
			if cows > 0 {
				fmt.Fprintf(out, "Коров: %d\n", cows)
			}
		}
	}
}
