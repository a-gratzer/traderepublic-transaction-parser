package writer

import (
	"github.com/a-gratzer/traderepublic-transaction-parser/internal/domain"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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

func (w *CSVWriter) MustWrite([]domain.MonthlyTransaction) {

	file := w.MustCreateFile()
	defer file.Close()

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
