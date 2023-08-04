package logfile

import (
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	file   *os.File
	errLog *log.Logger
)

func InitLog(fileName string) error {
	var err error
	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	errLog = log.New(file, "", log.Ldate|log.Ltime)
	errLog.SetOutput(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1,  // megabytes after which new file is created
		MaxBackups: 3,  // number of backups
		MaxAge:     15, // days
	})

	return nil
}

func DeInitLog() {
	file.Close()
}

func Printf(format string, v ...interface{}) {
	errLog.Printf(format, v...)
}

func Println(v ...interface{}) {
	errLog.Println(v...)
}

func Fatal(v ...interface{}) {
	errLog.Fatal(v...)
}
