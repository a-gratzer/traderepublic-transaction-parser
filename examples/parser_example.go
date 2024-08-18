package main

import (
	parser2 "github.com/a-gratzer/traderepublic-transaction-parser/internal/parser"
)

const (
	TEST_FILE = "examples/transactions-traderepublic.txt"
)

func main() {

	parser := parser2.NewTradeRepublicTransactionParser()
	transactions, _ := parser.MustParse(TEST_FILE)
	print(transactions)

}
