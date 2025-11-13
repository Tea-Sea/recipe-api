package logger

import (
	"log"
	"os"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

func NewAppLogger(prefix string, debug bool) *log.Logger {
	flags := log.LstdFlags
	if debug {
		flags |= log.LstdFlags | log.Lshortfile
	}
	return log.New(os.Stdout, prefix+" ", flags)
}

func NewGormLogger() gormlogger.Interface {
	return gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags|log.Lshortfile),
		gormlogger.Config{
			SlowThreshold: time.Second, // log slow SQL queries
			LogLevel:      gormlogger.Info,
			Colorful:      true,
		},
	)
}
