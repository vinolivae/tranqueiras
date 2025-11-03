package main

// Fontes
// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

func main() {
	// --- Dados de Teste ---
	key := []byte("minha-chave-secreta-muito-longa")
	message := []byte("Estes sao os dados que quero autenticar")
	password := []byte("Minh@SenhaForte!123")

	// O Argon2 requer um "salt" (dado aleatório).
	// Para um benchmark, podemos usar um estático.
	salt := []byte("somesalt16bytes") // 16 bytes

	fmt.Println("Iniciando benchmarks...\n")

	// ==========================================================
	// CENÁRIO 1: Benchmark HMAC-SHA256 (Foco em Throughput)
	// ==========================================================
	fmt.Println("--- Benchmark HMAC-SHA256 ---")
	fmt.Println("Medindo quantas operações em 1 segundo...")

	operations := 0
	startTime := time.Now()

	// Loop por 1 segundo
	for time.Since(startTime) < time.Second {
		// Esta é a operação de verificação
		h := hmac.New(sha256.New, key)
		h.Write(message)
		h.Sum(nil)
		operations++
	}

	totalTime := time.Since(startTime)

	fmt.Printf("Tempo total: %s\n", totalTime)
	fmt.Printf("Operações concluídas: %d\n", operations)
	opsPerSec := float64(operations) / totalTime.Seconds()
	fmt.Printf("Resultado: %.0f operações por segundo\n", opsPerSec)
	fmt.Println(strings.Repeat("-", 50))

	// ==========================================================
	// CENÁRIO 2: Benchmark Argon2id (Foco em Latência/Custo)
	// ==========================================================
	fmt.Println("\n--- Benchmark Argon2id ---")
	fmt.Println("Medindo o custo de UMA operação de hash e UMA de verificação...")

	// Parâmetros de Custo Minimo (Baseado nas recomendações OWASP)
	var timeCost uint32 = 5      // Iterações
	var memoryCost uint32 = 7168 // 7 MiB em KiB
	var parallelism uint8 = 5    // Threads
	var keyLen uint32 = 32       // Tamanho da saída (32 bytes)

	// 1. Medir o tempo para CRIAR O HASH
	startHash := time.Now()
	hash1 := argon2.IDKey(password, salt, timeCost, memoryCost, parallelism, keyLen)
	hashTime := time.Since(startHash)
	fmt.Printf("Tempo para criar UM hash: %s (%.2f ms)\n", hashTime, float64(hashTime.Nanoseconds())/1_000_000.0)

	// 2. Medir o tempo para VERIFICAR O HASH
	// A "verificação" em Go puro envolve rodar o hash novamente com os
	// mesmos parâmetros (salt, etc.) e comparar o resultado em tempo constante.
	startVerify := time.Now()
	hash2 := argon2.IDKey(password, salt, timeCost, memoryCost, parallelism, keyLen)
	match := hmac.Equal(hash1, hash2)
	verifyTime := time.Since(startVerify)

	fmt.Printf("Tempo para verificar UM hash: %s (%.2f ms)\n", verifyTime, float64(verifyTime.Nanoseconds())/1_000_000.0)
	fmt.Printf("Resultado da verificação: %t\n", match)
	fmt.Println(strings.Repeat("-", 50))
}
