package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var LOGGER *logrus.Logger

func New(viper *viper.Viper) *os.File {
	outputFilePath := viper.GetString("log.outputFile")
	logLevel := viper.GetInt32("log.level")

	dir := filepath.Dir(outputFilePath)
	// check if the file output exists
	// if doesn't exist, create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal("failed to create file:", err)
		}
	}

	file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("failed to open file:", err)
	}

	log := logrus.New()

	log.SetLevel(logrus.Level(logLevel))
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	log.SetOutput(file)

	LOGGER = log

	return file
}

func Get() *logrus.Logger {
	if LOGGER == nil {
		log.Fatal("logger is not initialized")
	}
	return LOGGER
}
