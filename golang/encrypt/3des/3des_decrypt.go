package tripleDES

import (
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

// KeyProvider simula a obtenção de chaves como a função `get_key` em Elixir.
// Em uma aplicação real, isso buscaria chaves de um local seguro.
func getKey(keyReference, network string) (string, error) {
	switch keyReference {
	case "pismo":
		if network == "MASTERCARD" {
			// Chave de exemplo para pismo_mastercard_pinblock_encryption
			return "0123456789ABCDEFFEDCBA9876543210", nil
		}
		// Chave de exemplo para pismo_pinblock_encryption
		return "FEDCBA98765432100123456789ABCDEF", nil
	case "orbitall":
		// Chave de exemplo para pinblock_encryption
		return "11112222333344445555666677778888", nil
	case "tutuka":
		// Chave de exemplo para tutuka_pinblock_encryption
		return "88887777666655554444333322221111", nil
	default:
		return "", fmt.Errorf("referência de chave desconhecida: %s", keyReference)
	}
}

// decrypt3DESECB decifra dados usando 3DES com uma chave de 16 bytes (opção de chaveamento 2) no modo ECB.
// Isso é feito criando uma chave de 24 bytes (k1+k2+k1) para compatibilidade com a biblioteca padrão do Go.
func decrypt3DESECB(key, ciphertext []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("o tamanho da chave deve ser de 16 bytes, mas foi de %d", len(key))
	}
	if len(ciphertext)%des.BlockSize != 0 {
		return nil, fmt.Errorf("o tamanho do ciphertext deve ser um múltiplo de %d", des.BlockSize)
	}

	// Constrói uma chave de 24 bytes a partir da chave de 16 bytes (k1+k2 -> k1+k2+k1)
	tripleDESKey := make([]byte, 24)
	copy(tripleDESKey, key[:16])
	copy(tripleDESKey[16:], key[:8])

	block, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += des.BlockSize {
		block.Decrypt(decrypted[i:i+des.BlockSize], ciphertext[i:i+des.BlockSize])
	}

	return decrypted, nil
}

// extractPinFromPinblock extrai o PIN do pinblock decifrado com base no formato ISO.
func extractPinFromPinblock(decryptedPinblock []byte, pan, isoFormat string) (string, error) {
	switch isoFormat {
	case "3":
		if len(pan) < 15 {
			return "", errors.New("PAN inválido para o formato ISO 3")
		}
		// Prepara o bloco do PAN: "0000" + 12 dígitos do PAN
		panPartHex := "0000" + pan[3:15]
		panBlock, err := hex.DecodeString(panPartHex)
		if err != nil {
			return "", fmt.Errorf("falha ao decodificar o bloco do PAN: %w", err)
		}

		// XOR do pinblock decifrado com o bloco do PAN
		pinPart := make([]byte, len(decryptedPinblock))
		for i := range decryptedPinblock {
			pinPart[i] = decryptedPinblock[i] ^ panBlock[i]
		}

		pinPartHex := hex.EncodeToString(pinPart)
		lenDigit := string(pinPartHex[1])
		length, err := strconv.ParseInt(lenDigit, 16, 64)
		if err != nil {
			return "", fmt.Errorf("falha ao analisar o comprimento do PIN: %w", err)
		}
		return pinPartHex[2 : 2+length], nil

	case "2":
		encodedPinblock := hex.EncodeToString(decryptedPinblock)
		lenDigit := string(encodedPinblock[0])
		length, err := strconv.ParseInt(lenDigit, 16, 64)
		if err != nil {
			return "", fmt.Errorf("falha ao analisar o comprimento do PIN: %w", err)
		}
		return encodedPinblock[1 : 1+length], nil

	case "1":
		encodedPinblock := hex.EncodeToString(decryptedPinblock)
		if string(encodedPinblock[0]) != "1" {
			return "", errors.New("formato ISO 1 inválido: não começa com '1'")
		}
		lenDigit := string(encodedPinblock[1])
		length, err := strconv.ParseInt(lenDigit, 16, 64)
		if err != nil {
			return "", fmt.Errorf("falha ao analisar o comprimento do PIN: %w", err)
		}
		return encodedPinblock[2 : 2+length], nil

	default:
		return "", fmt.Errorf("formato ISO não suportado: %s", isoFormat)
	}
}

// decryptBlock orquestra o processo de decifragem.
func DecryptBlock(blockHex, pan, keyReference, network, isoFormat string) (string, error) {
	keyHex, err := getKey(keyReference, network)
	if err != nil {
		return "", err
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", fmt.Errorf("falha ao decodificar a chave: %w", err)
	}

	block, err := hex.DecodeString(blockHex)
	if err != nil {
		return "", fmt.Errorf("falha ao decodificar o pinblock: %w", err)
	}

	decryptedPinblock, err := decrypt3DESECB(key, block)
	if err != nil {
		return "", fmt.Errorf("falha ao decifrar o pinblock: %w", err)
	}

	return extractPinFromPinblock(decryptedPinblock, pan, isoFormat)
}
