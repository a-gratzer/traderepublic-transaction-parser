package parser

import (
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMustParse_this_month_token(t *testing.T) {

	parser := NewTradeRepublicTransactionParser(logger.GetZapLogger(false))
	data, err := parser.MustParse("testdata/parser_test_input_this_month.txt")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, time.Now().Year(), data[0].Year)
	assert.Equal(t, time.Now().Month(), data[0].Month)

}

func TestMustParse_month_tokens(t *testing.T) {

	parser := NewTradeRepublicTransactionParser(logger.GetZapLogger(false))
	data, err := parser.MustParse("testdata/parser_test_input_month_tokens.txt")
	assert.NoError(t, err)
	assert.Equal(t, 12, len(data))

	for _, d := range data {
		assert.Equal(t, time.Now().Year(), d.Year, "All have current year - no year set on token")
	}

	assert.Equal(t, time.January, data[0].Month, "First month in test-data is January")
	assert.Equal(t, time.February, data[1].Month, "2nd month in test-data is February")
	assert.Equal(t, time.March, data[2].Month, "3rd month in test-data is March")
	assert.Equal(t, time.December, data[11].Month, "12th month in test-data is December")
}

func TestMustParse_month_year_tokens(t *testing.T) {

	parser := NewTradeRepublicTransactionParser(logger.GetZapLogger(false))
	data, err := parser.MustParse("testdata/parser_test_input_month_year_tokens.txt")
	assert.NoError(t, err)
	assert.Equal(t, 12, len(data))

	for _, d := range data {
		assert.Equal(t, 2023, d.Year, "All have year 2023 because 2023 is set on token")
	}

	assert.Equal(t, time.January, data[0].Month, "First month in test-data is January")
	assert.Equal(t, time.February, data[1].Month, "2nd month in test-data is February")
	assert.Equal(t, time.March, data[2].Month, "3rd month in test-data is March")
	assert.Equal(t, time.December, data[11].Month, "12th month in test-data is December")
}

func TestMustParse_limit_buy_tokens(t *testing.T) {

	parser := NewTradeRepublicTransactionParser(logger.GetZapLogger(false))
	data, err := parser.MustParse("testdata/parser_test_input_limit_buy.txt")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(data))

	assert.Equal(t, 3, len(data[0].Transactions), "3 buy orders expected")

	assert.Equal(t, time.January, data[0].Transactions[0].Date.Month(), "01/01 is January")
	assert.Equal(t, time.January, data[0].Transactions[1].Date.Month(), "02/01 is January")
	assert.Equal(t, time.January, data[0].Transactions[2].Date.Month(), "03/01 is January")

	assert.Equal(t, 1, data[0].Transactions[0].Date.Day(), "01/01 is 1")
	assert.Equal(t, 2, data[0].Transactions[1].Date.Day(), "02/01 is 2")
	assert.Equal(t, 3, data[0].Transactions[2].Date.Day(), "03/01 is 3")

}
