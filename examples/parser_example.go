package main

import (
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/config"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/logger"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/parser"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/writer"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	config.InitDefaultViperConfig(getConfigFilePath())
	log = logger.GetZapLogger(false)
}

func main() {

	tradeRepParser := parser.NewTradeRepublicTransactionParser(log)
	result, err := tradeRepParser.MustParse(getInputFilePath())
	if err != nil {
		log.Error(err.Error())
	}
	writer.NewCSVWriter(log, getOutputFilePath()).MustWrite(result)
	print(result)

}

func getConfigFilePath() string {
	return "./config/config.yaml"
}

func getOutputFilePath() string {
	return viper.GetViper().GetString(parser.CONFIG_PARSER_OUTPUT_FILE)
}

func getInputFilePath() string {
	return viper.GetViper().GetString(parser.CONFIG_PARSER_INPUT_FILE)
}
