package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"db-rodrigo-m/internal/storage"
)

func main() {
	db, err := storage.NewDB("wal.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao abrir banco:", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("Mini DB")
	fmt.Println("Comandos:")
	fmt.Println("  SET chave valor")
	fmt.Println("  GET chave")
	fmt.Println("  DEL chave")
	fmt.Println("  EXIT")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("db > ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "SET":
			if len(parts) < 3 {
				fmt.Println("uso: SET chave valor")
				continue
			}

			key := parts[1]
			value := strings.Join(parts[2:], " ")

			if err := db.Set(key, value); err != nil {
				fmt.Println("erro:", err)
				continue
			}

			fmt.Println("OK")

		case "GET":
			if len(parts) != 2 {
				fmt.Println("uso: GET chave")
				continue
			}

			key := parts[1]
			value, ok := db.Get(key)
			if !ok {
				fmt.Println("not found")
				continue
			}

			fmt.Println(value)

		case "DEL":
			if len(parts) != 2 {
				fmt.Println("uso: DEL chave")
				continue
			}

			key := parts[1]

			if err := db.Delete(key); err != nil {
				fmt.Println("erro:", err)
				continue
			}

			fmt.Println("OK")

		case "EXIT", "QUIT":
			fmt.Println("bye")
			return

		default:
			fmt.Println("comando desconhecido")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "erro lendo entrada:", err)
	}
}