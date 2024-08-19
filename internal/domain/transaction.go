package domain

import "time"

var MonthMap = map[string]time.Month{
	"January":   time.January,
	"February":  time.February,
	"March":     time.March,
	"April":     time.April,
	"May":       time.May,
	"June":      time.June,
	"July":      time.July,
	"August":    time.August,
	"September": time.September,
	"October":   time.October,
	"November":  time.November,
	"December":  time.December,
}

type MonthlyTransaction struct {
	Year         int
	Month        time.Month
	Transactions []Transaction
}

type Transaction struct {
	Date time.Time
	Tag  string
	Type string

	Amount Amount
	Raw    []string
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

type Amount struct {
	Currency string
	Prefix   string
	AbsValue float64
}
