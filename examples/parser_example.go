package main

import (
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/config"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/logger"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/parser"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	config.InitDefaultViperConfig("./config/config.yaml")
	log = logger.GetZapLogger(false)
}

func main() {

	tradeRepParser := parser.NewTradeRepublicTransactionParser(log)
	result, _ := tradeRepParser.MustParse(viper.GetViper().GetString(parser.CONFIG_PARSER_INPUT_FILE))
	print(result)

}
