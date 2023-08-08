package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	LogDirpath string
	LogFile    *os.File
}

func NewLogger(dirpath, usecase string) *Logger {
	logger := &Logger{LogDirpath: dirpath}
	logger.createLogfile(usecase)
	return logger
}

func (lg *Logger) createLogfile(prefix string) error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(lg.LogDirpath, os.ModePerm); err != nil {
		log.Fatal("Error creating logs directory:", err)
		return err
	}

	// Get current date and time
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	// create filename suffixed with current date and time
	logFileName := fmt.Sprintf("%s_%s", prefix, currentTime)
	// Create a log file in the logs directory with the filename
	file, err := os.Create(filepath.Join(lg.LogDirpath, logFileName))
	if err != nil {
		log.Fatal("Error creating log file:", err)
		return err
	}
	file.Close()
	lg.LogFile = file
	return nil
}

func (lg *Logger) Log(level, details string) {
	logger := log.New(lg.LogFile,
		fmt.Sprintf("%s: ", level), log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(details)
}
