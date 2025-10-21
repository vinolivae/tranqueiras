package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2Params struct {
	m    uint32
	t    uint32
	p    uint8
	salt string
}

func main() {
	password := "123456"
	key := make([]byte, 32)

	hashedPassword, err := hashPassword(password, key)
	if err != nil {
		fmt.Println("Error encoding password:", err)
		return
	}
	fmt.Println("old:", hashedPassword)

	isValid, err := verifyPassword(password, hashedPassword, key)
	if err != nil {
		fmt.Println("Error verifying password:", err)
		return
	}
	fmt.Println("Is Password valid?", isValid)
}

func hashPassword(password string, key []byte) (string, error) {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(password))
	keyedHash := mac.Sum(nil)

	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	t := uint32(3)
	m := uint32(10 * 1024)
	p := uint8(1)
	keyLen := uint32(32)

	hash := argon2.IDKey(keyedHash, salt, t, m, p, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, m, t, p, b64Salt, b64Hash), nil
}

func verifyPassword(password, expectedHash string, key []byte) (bool, error) {
	params, _ := parseArgon2Params(expectedHash)

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(password))
	keyedHash := mac.Sum(nil)

	hash := argon2.IDKey(keyedHash, []byte(params.salt), params.t, params.m, params.p, 32)
	b64Salt := base64.RawStdEncoding.EncodeToString([]byte(params.salt))
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	hashedPassword := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.m, params.t, params.p, b64Salt, b64Hash)
	fmt.Println("new:", hashedPassword)
	return subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(expectedHash)) == 1, nil
}

func parseArgon2Params(hashSaved string) (params argon2Params, err error) {
	parts := strings.Split(hashSaved, "$")
	if len(parts) != 6 {
		return argon2Params{}, errors.New("invalid hash format")
	}

	paramSection := parts[3]
	for _, v := range strings.Split(paramSection, ",") {
		values := strings.SplitN(v, "=", 2)
		if len(values) != 2 {
			return argon2Params{}, errors.New("invalid parameter section in hash")
		}

		switch values[0] {
		case "m":
			val, err := strconv.Atoi(values[1])
			if err != nil {
				return argon2Params{}, errors.New("parsing memory parameter")
			}
			params.m = uint32(val)
		case "t":
			val, err := strconv.Atoi(values[1])
			if err != nil {
				return argon2Params{}, errors.New("parsing time parameter")
			}
			params.t = uint32(val)
		case "p":
			val, err := strconv.Atoi(values[1])
			if err != nil {
				return argon2Params{}, errors.New("parsing parallelism parameter")
			}
			params.p = uint8(val)
		default:
			return argon2Params{}, fmt.Errorf("unknown parameter in hash: %s", values[0])
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return argon2Params{}, fmt.Errorf("decoding base64 salt: %w", err)
	}
	params.salt = string(salt)

	return params, nil
}
