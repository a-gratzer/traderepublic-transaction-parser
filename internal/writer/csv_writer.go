package writer

import (
	"encoding/csv"
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/domain"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strconv"
)

type CSVWriter struct {
	filePath string
	logger   *zap.Logger
}

func NewCSVWriter(logger *zap.Logger, filePath string) *CSVWriter {
	return &CSVWriter{
		filePath: filePath,
		logger:   logger,
	}
}

func (w *CSVWriter) MustWrite(transactions []domain.MonthlyTransaction) {

	file := w.MustCreateFile()
	defer file.Close()

	writer := csv.NewWriter(file)
	w.writeHeader(writer)
	w.writeTransactions(writer, transactions)

	writer.Flush()

}

func (w *CSVWriter) MustCreateFile() *os.File {
	dir := filepath.Dir(w.filePath)
	parent := filepath.Base(dir)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		os.MkdirAll(parent, 0700)
	}

	file, err := os.Create(w.filePath)
	if err != nil {
		w.logger.Error("Failed to create file", zap.String("path", w.filePath), zap.Error(err))
	}

	return file
}

func (w *CSVWriter) writeHeader(writer *csv.Writer) {
	err := writer.Write([]string{"Date", "Tag", "Type", "Currency", "Amount_Prefix", "Amount_Absolute"})
	if err != nil {
		w.logger.Error("Failed to write header", zap.Error(err))
	}
}

func (w *CSVWriter) writeTransactions(writer *csv.Writer, transactions []domain.MonthlyTransaction) {
	for _, t := range transactions {
		for _, mt := range t.Transactions {
			writer.Write([]string{
				mt.Date.Format("2006-01-02"),
				mt.Tag,
				mt.Type,
				mt.Amount.Currency,
				mt.Amount.Prefix,
				strconv.FormatFloat(mt.Amount.AbsValue, 'f', -1, 64),
			})

		}
	}
}
