package domain

import "time"

type MonthlyTransaction struct {
	Year         int
	Month        int
	Transactions []Transaction
}

func NewMonthlyTransaction(date time.Time) *MonthlyTransaction {
	return &MonthlyTransaction{
		Year:  date.Year(),
		Month: int(date.Month()),
	}
}

type Transaction struct {
	Date time.Time
	Tag  string
	Type string

	Amount Amount
	Raw    []string
}

type Amount struct {
	Currency string
	Value    float64
}
