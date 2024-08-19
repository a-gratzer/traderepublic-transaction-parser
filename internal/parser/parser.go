package parser

import (
	"bufio"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/domain"
	"os"
	"regexp"
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

	monthly := make([]domain.MonthlyTransaction, 0)
	var currentMonth *domain.MonthlyTransaction = nil
	var currentTransaction *domain.Transaction = domain.NewTransaction()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()

		if line == "" {
			continue
		}

		if p.isYearMonthToken(line) {
			if currentMonth == nil {
				currentMonth = p.mustGetMonthlyTransaction(line)
			} else {
				monthly = append(monthly, *currentMonth)
				currentMonth = p.mustGetMonthlyTransaction(line)
			}

		} else {
			currentTransaction.Raw = append(currentTransaction.Raw, line)
			if p.isPriceToken(line) {
				p.mustParseTransaction(currentMonth.Year, currentMonth.Month, currentTransaction)
				currentMonth.Transactions = append(currentMonth.Transactions, *currentTransaction)
				currentTransaction = domain.NewTransaction()
			}
		}
	}

	return monthly, nil
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

func (p *TradeRepublicTransactionParser) isDayAndTypeToken(line string) bool {
	DAY_TYPE_PATTERN := `^(\d{2}\/\d{2})(.*?)$`

	re := regexp.MustCompile(DAY_TYPE_PATTERN)

	return re.MatchString(line)
}

func (p *TradeRepublicTransactionParser) isPriceToken(line string) bool {
	AMOUNT_PATTERN := `^(\+|-)?(€)([0-9.,]+)$`
	re := regexp.MustCompile(AMOUNT_PATTERN)
	return re.MatchString(line)
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
			date = time.Date(year, month, 0, 0, 0, 0, 0, time.UTC)
		}

	}

	return &domain.MonthlyTransaction{
		Year:  date.Year(),
		Month: date.Month(),
	}
}

func (p *TradeRepublicTransactionParser) mustParseTransaction(year int, month time.Month, transaction *domain.Transaction) {
	transaction.Tag = transaction.Raw[0]

}
