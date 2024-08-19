package main

import (
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/config"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/logger"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/parser"
	"github.com/spf13/viper"
)

const (
	TEST_FILE = "examples/transactions-traderepublic.txt"
)

func main() {

	config.InitDefaultViperConfig("./config/config.yaml")
	user := viper.GetViper().GetString("facts.user")
	print(user)
	logger := logger.GetZapLogger(false)

	parser.NewTradeRepublicTransactionParser(logger)
	//if transactions, err := tradeRepParser.MustParse(TEST_FILE); err != nil {
	//	for _, monthlyTransactions := range transactions {
	//		print(string(monthlyTransactions))
	//	}
	//}

}
