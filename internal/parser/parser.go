package parser

import (
	"bufio"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/domain"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	AMOUNT_PATTERN   = `^(\+|-)?({SYMBOL})([0-9{THOUSAND}{DECIMAL}]+)$`
	DAY_TYPE_PATTERN = `^(\d{2}\/\d{2})(.*?)$`
)

type TradeRepublicTransactionParser struct {
	logger    *zap.Logger
	amountXp  *regexp.Regexp
	dayTypeXp *regexp.Regexp
}

func NewTradeRepublicTransactionParser(logger *zap.Logger) *TradeRepublicTransactionParser {
	parser := &TradeRepublicTransactionParser{
		logger: logger,
	}
	parser.applyConfig()
	return parser
}

func (p *TradeRepublicTransactionParser) applyConfig() {
	viper.SetDefault(CONFIG_PARSER_TRANSACTION_AMOUNT_SYMBOL, "€")
	viper.SetDefault(CONFIG_PARSER_TRANSACTION_AMOUNT_SEPARATOR_THOUSAND, ",")
	viper.SetDefault(CONFIG_PARSER_TRANSACTION_AMOUNT_SEPARATOR_DECIMAL, ".")
	amountPattern := strings.Replace(AMOUNT_PATTERN, "{SYMBOL}", viper.GetViper().GetString(CONFIG_PARSER_TRANSACTION_AMOUNT_SYMBOL), 1)
	amountPattern = strings.Replace(amountPattern, "{THOUSAND}", viper.GetViper().GetString(CONFIG_PARSER_TRANSACTION_AMOUNT_SEPARATOR_THOUSAND), 1)
	amountPattern = strings.Replace(amountPattern, "{DECIMAL}", viper.GetViper().GetString(CONFIG_PARSER_TRANSACTION_AMOUNT_SEPARATOR_DECIMAL), 1)

	p.amountXp = regexp.MustCompile(amountPattern)
	p.dayTypeXp = regexp.MustCompile(DAY_TYPE_PATTERN)
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
				currentMonth = p.mustGetMonthlyTransaction(line)
			}
			monthly = append(monthly, *currentMonth)
		} else {
			currentTransaction.Raw = append(currentTransaction.Raw, line)
			if p.isPriceToken(line) {
				// last entry from one transaction
				p.mustParseTransaction(currentMonth.Year, currentTransaction)
				last := &monthly[len(monthly)-1]
				last.Transactions = append(last.Transactions, *currentTransaction)
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
	return p.dayTypeXp.MatchString(line)
}

func (p *TradeRepublicTransactionParser) isPriceToken(line string) bool {
	return p.amountXp.MatchString(line)
}

func (p *TradeRepublicTransactionParser) mustGetMonthlyTransaction(line string) *domain.MonthlyTransaction {

	date := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	if strings.Contains(line, "This month") {
		date = time.Now().UTC()
	} else {
		parts := strings.Split(line, " ")
		switch len(parts) {
		case 0:
			p.logger.Error("Unable to parse Month/Year token")
		case 1:
			month, _ := domain.MonthMap[parts[0]]
			date = time.Date(time.Now().Year(), month, 1, 0, 0, 0, 0, time.UTC)
		case 2:
			month, _ := domain.MonthMap[parts[0]]
			year, _ := strconv.Atoi(parts[1])
			date = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		}
	}

	return &domain.MonthlyTransaction{
		Year:  date.Year(),
		Month: date.Month(),
	}
}

func (p *TradeRepublicTransactionParser) mustParseTransaction(year int, transaction *domain.Transaction) {

	if len(transaction.Raw) != 4 {
		p.logger.Error("Unable to parse transaction.",
			zap.Int("Expected tokens", 4),
			zap.Int("Observed tokens", len(transaction.Raw)),
		)
	}

	// #############################
	// Tag
	transaction.Tag = transaction.Raw[0]

	// #############################
	// Day/Type
	matchDayType := p.dayTypeXp.FindStringSubmatch(transaction.Raw[2])
	if len(matchDayType) != 3 {
		p.logger.Error("Unable to parse day/type token",
			zap.String("token", transaction.Raw[2]),
			zap.Strings("match", matchDayType),
		)
	}

	dayMonth := strings.Split(matchDayType[1], "/")
	day, _ := strconv.Atoi(dayMonth[0])
	month, _ := strconv.Atoi(dayMonth[1])
	transaction.Date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	transaction.Type = strings.Replace(matchDayType[2], " - ", "", 1)

	// #############################
	// Amount
	matchAmount := p.amountXp.FindStringSubmatch(transaction.Raw[3])
	if len(matchAmount) != 4 {
		p.logger.Error("Unable to parse amount token ",
			zap.String("token", transaction.Raw[3]),
			zap.Strings("match", matchAmount),
		)
	}
	transaction.Amount.Currency = matchAmount[2]

	val, err := strconv.ParseFloat(strings.Replace(matchAmount[3], ",", "", 5), 64)
	if err != nil {
		p.logger.Error("Unable to parse amount ",
			zap.String("token", transaction.Raw[3]),
		)
	}

	transaction.Amount.Prefix = matchAmount[1]
	if matchAmount[1] == "" {
		transaction.Amount.Prefix = "-"
	}

	transaction.Amount.AbsValue = val

}
