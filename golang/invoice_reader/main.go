package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Invoice struct {
	PurchaseDate   string
	Estabilishment string
	Bearer         string
	Amount         float64
	Portion        string
}

func main() {
	records := readFile()
	invoices := formatReadedFile(records)
	categorizedInvoice := categorizeInvoiceAmount(invoices)
	totalAmount := calculateTotalAmount(categorizedInvoice)
	writeFile(categorizedInvoice, totalAmount)
}

func readFile() [][]string {
	file, err := os.Open("invoice.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Define o delimitador como ponto e vírgula

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return records
}

func formatReadedFile(records [][]string) []Invoice {
	invoices := make([]Invoice, len(records))
	for k, v := range records {
		strWithoutPrefix, _ := strings.CutPrefix(v[3], "R$ ")
		strWithDot := strings.ReplaceAll(strWithoutPrefix, ",", ".")
		amount, err := strconv.ParseFloat(strWithDot, 64)

		if err != nil {
			panic(err)
		}

		invoice := Invoice{
			PurchaseDate:   v[0],
			Estabilishment: v[1],
			Bearer:         v[2],
			Amount:         amount,
			Portion:        v[4],
		}

		invoices[k] = invoice
	}

	return invoices
}

func categorizeInvoiceAmount(invoices []Invoice) map[string]float64 {
	categories := make(map[string]float64)
	for _, invoice := range invoices {
		category := possibleCategories(invoice.Estabilishment)
		categories[category] += invoice.Amount
	}

	for category, total := range categories {
		categories[category] = math.Floor(total*100) / 100
	}

	return categories
}

func calculateTotalAmount(categorizedInvoice map[string]float64) float64 {
	var totalAmount float64

	for _, total := range categorizedInvoice {
		totalAmount += total
	}

	return totalAmount
}

func possibleCategories(estabilishmentName string) string {
	categories := map[string][]string{
		"Restaurant": {"ifd", "ifood"},
		"Transport":  {"uber"},
		"Leisure":    {"spotify", "steam"},
	}

	for category, keywords := range categories {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(estabilishmentName), strings.ToLower(keyword)) {
				return category
			}
		}
	}

	return "Other"
}

func writeFile(categorizedInvoice map[string]float64, totalAmount float64) {
	file, err := os.Create("categorized_invoice.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	records := formatToWriteFile(categorizedInvoice, totalAmount)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}
}

func formatToWriteFile(categorizedInvoice map[string]float64, totalAmount float64) [][]string {
	records := [][]string{
		{"Category", "Total"},
	}

	for category, total := range categorizedInvoice {
		records = append(records, []string{category, fmt.Sprintf("%.2f", total)})
	}

	records = append(records, []string{"Total", fmt.Sprintf("%.2f", totalAmount)})

	return records
}
