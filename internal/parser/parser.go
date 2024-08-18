package parser

import (
	"bufio"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/domain"
	"os"
	"strconv"
	"strings"
	"time"
)

type TradeRepublicTransactionParser struct {
}

func NewTradeRepublicTransactionParser() *TradeRepublicTransactionParser {
	return &TradeRepublicTransactionParser{}
}

func (p *TradeRepublicTransactionParser) MustParse(filePath string) ([]domain.MonthlyTransaction, error) {

	file := p.mustOpenFile(filePath)
	defer file.Close()

	transactions := make([]domain.MonthlyTransaction, 0)

	var currentTransaction *domain.MonthlyTransaction = nil
	// Create a scanner
	scanner := bufio.NewScanner(file)

	// Read and print lines
	for scanner.Scan() {

		line := scanner.Text()

		if p.isYearMonthToken(line) {
			currentTransaction = p.mustGetMonthlyTransaction(line)
			transactions = append(transactions, *currentTransaction)
		} else {
			currentTransaction.Transactions = append(currentTransaction.Transactions)
		}
	}

	return transactions, nil
}

func (p *TradeRepublicTransactionParser) mustOpenFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

func (p *TradeRepublicTransactionParser) isYearMonthToken(line string) bool {

	if strings.Contains(line, "This month") {
		return true
	}

	for monthString := range domain.MonthMap {
		if strings.Contains(line, monthString) {
			return true
		}
	}

	return false
}

func (p *TradeRepublicTransactionParser) mustGetMonthlyTransaction(line string) *domain.MonthlyTransaction {

	date := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	if strings.Contains(line, "This month") {
		date = date.AddDate(time.Now().Year(), int(time.Now().Month()), 0)
	} else {
		parts := strings.Split(line, " ")

		switch len(parts) {
		case 0:
			panic("Unable to parse Month/Year token")
		case 1:
			month, _ := domain.MonthMap[parts[0]]
			date = time.Date(time.Now().Year(), month, 0, 0, 0, 0, 0, time.UTC)
		case 2:
			month, _ := domain.MonthMap[parts[0]]
			year, _ := strconv.Atoi(parts[1])
			date = time.Date(year, month, 0, 0, 0, 0, 1, time.UTC)
		}

	}

	return &domain.MonthlyTransaction{
		Year:  date.Year(),
		Month: date.Month(),
	}
}
