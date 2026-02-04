package main

import (
	"log"

	tripleDES "github.com/vinolivae/tranqueiras/tree/main/golang/encrypt/3des"
)

func main() {
	// Dados de exemplo baseados na lógica do Elixir.
	// Substitua por seus dados de teste reais.
	input := map[string]string{
		"pinblock":      "11A554BEA9A455F1",
		"pan":           "1234567890123456",
		"key_reference": "pismo",
		"network":       "MASTERCARD",
		"iso_format":    "3",
	}

	log.Printf("Decifrando PIN para PAN: %s", input["pan"])

	pin, err := tripleDES.DecryptBlock(
		input["pinblock"],
		input["pan"],
		input["key_reference"],
		input["network"],
		input["iso_format"],
	)

	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	log.Printf("PIN decifrado com sucesso: %s", pin)
}
