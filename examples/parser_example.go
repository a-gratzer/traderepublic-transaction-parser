package main

import (
	"bufio"
	"fmt"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/domain"
	"os"
	"time"
)

const (
	TEST_FILE = "examples/transactions-traderepublic.txt"
)

func main() {

	// Open the file
	file := MustOpenFile(TEST_FILE)
	defer file.Close()

	monthly := domain.NewMonthlyTransaction(time.Now())

	// Create a scanner
	scanner := bufio.NewScanner(file)

	// Read and print lines
	for scanner.Scan() {

		line := scanner.Text()

		fmt.Println(line)

	}

	// Check for errors
	fmt.Println("#####################")
	fmt.Println("Errors:")
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func MustOpenFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return file
}
