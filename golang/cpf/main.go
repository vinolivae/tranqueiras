package main

import (
	"errors"
	"fmt"
	"strconv"
)

func main() {
	CPF := "34485861023"
	// CPF := "17674584722"
	// CPF := "176745847"

	if err := validateLength(CPF); err != nil {
		panic(err)

	}

	if err := validateFirstVerificationDigit(CPF); err != nil {
		panic(err)
	}

	if err := validateSecondVerificationDigit(CPF); err != nil {
		panic(err)
	}

	fmt.Println("CPF Válido!")
	fmt.Println("Esse CPF pertence ao municipio de", discoverRegion(string(CPF[7])))
}

func validateLength(cpf string) error {
	if len(cpf) != 11 {
		return errors.New("Invalid CPF Length")
	}
	return nil
}

func validateFirstVerificationDigit(cpf string) error {
	firstVerificationDigit, _ := strconv.Atoi(string(cpf[9]))
	var sumOfFirstNineDigits int
	digitLenSlice := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	for i := 0; i < 9; i++ {
		cpfDigit, _ := strconv.Atoi(string(cpf[i]))
		sumOfFirstNineDigits += cpfDigit * digitLenSlice[i]
	}

	divOfSumOfFirstNineDigits := sumOfFirstNineDigits / 11

	result := sumOfFirstNineDigits - (11 * divOfSumOfFirstNineDigits)

	verifyNumberZeroSlice := []int{0, 1}
	verifyOtherNumberSlice := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, v := range verifyNumberZeroSlice {
		if result == v && firstVerificationDigit == 0 {
			return nil
		}
	}

	for _, v := range verifyOtherNumberSlice {
		expectedFirstDigit := 11 - result

		if result == v && firstVerificationDigit == expectedFirstDigit {
			return nil
		}
	}

	return errors.New("Invalid First Verification Digit")
}

func validateSecondVerificationDigit(cpf string) error {
	secondVerificationDigit, _ := strconv.Atoi(string(cpf[10]))
	var sumOfFirstTenDigits int
	digitLenSlice := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
	for i := 0; i < 10; i++ {
		cpfDigit, _ := strconv.Atoi(string(cpf[i]))
		sumOfFirstTenDigits += cpfDigit * digitLenSlice[i]
	}

	divOfSumOfFirstTenDigits := sumOfFirstTenDigits / 11

	result := sumOfFirstTenDigits - (11 * divOfSumOfFirstTenDigits)

	verifyNumberZeroSlice := []int{0, 1}
	verifyOtherNumberSlice := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, v := range verifyNumberZeroSlice {
		if result == v && secondVerificationDigit == 0 {
			return nil
		}
	}

	for _, v := range verifyOtherNumberSlice {
		expectedFirstDigit := 11 - result

		if result == v && secondVerificationDigit == expectedFirstDigit {
			return nil
		}
	}

	return errors.New("Invalid Second Verification Digit")
}

func discoverRegion(eighthDigit string) string {
	municipality := map[string]string{
		"0": "Rio Grande do Sul",
		"1": "Distrito Federal, Goiás, Mato Grosso, Mato Grosso do Sul ou Tocantins",
		"2": "Amazonas, Pará, Roraima, Amapá, Acre ou Rondônia",
		"3": "Ceará, Maranhão ou Piauí",
		"4": "Paraíba, Pernambuco, Alagoas ou Rio Grande do Norte",
		"5": "Bahia ou Sergipe",
		"6": "Minas Gerais",
		"7": "Rio de Janeiro ou Espírito Santo",
		"8": "São Paulo",
		"9": "Paraná e Santa Catarina",
	}

	return municipality[eighthDigit]
}
