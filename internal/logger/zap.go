package logger

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLogger *zap.Logger = nil

func GetZapLogger(forceDebug bool) *zap.Logger {
	if zapLogger == nil {
		SetupZapLogger(forceDebug)
	}
	return zapLogger
}

func SetupZapLogger(forceDebug bool) {
	var loggerConfig zap.Config
	var cmdDebug = viper.GetViper().GetBool(CONFIG_LOGGER_DEBUG)
	if cmdDebug || forceDebug {
		loggerConfig = zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	var outputPaths = viper.GetViper().GetStringSlice(CONFIG_LOGGER_OUTPUT_LIST)
	if outputPaths != nil && len(outputPaths) > 0 {
		loggerConfig.OutputPaths = outputPaths
	}

	var logLevel = viper.GetViper().GetString(CONFIG_LOGGER_LOG_LEVEL)
	if logLevel == "INFO" {
		loggerConfig.Level.SetLevel(zap.InfoLevel)
	} else if logLevel == "DEBUG" {
		loggerConfig.Level.SetLevel(zap.DebugLevel)
	} else if logLevel == "WARN" {
		loggerConfig.Level.SetLevel(zap.WarnLevel)
	} else if logLevel == "ERROR" {
		loggerConfig.Level.SetLevel(zap.ErrorLevel)
	} else if logLevel == "FATAL" {
		loggerConfig.Level.SetLevel(zap.FatalLevel)
	} else if logLevel == "PANIC" {
		loggerConfig.Level.SetLevel(zap.PanicLevel)
	} else {
		loggerConfig.Level.SetLevel(zap.InfoLevel)
	}

	l, err := loggerConfig.Build(
		zap.AddStacktrace(zapcore.FatalLevel), // only print stack traces for fatal errors
	)
	if err != nil {
		fmt.Errorf("failed to set up logging: %w", err)
		return
	}

	zapLogger = l

	if cmdDebug || forceDebug {
		zapLogger.Info("###################################")
		zapLogger.Info("Logger config", zap.String("logLevel", loggerConfig.Level.String()))
		zapLogger.Debug("Logger config", zap.String("log-message", "Debug message"))
		zapLogger.Info("Logger config", zap.String("log-message", "Info message"))
		zapLogger.Warn("Logger config", zap.String("log-message", "Warn message"))
		zapLogger.Error("Logger config", zap.String("log-message", "Error message"))
		zapLogger.Info("###################################")
		//logger.Fatal("Logger initialized")
		//logger.Panic("Logger initialized")
	}
}
