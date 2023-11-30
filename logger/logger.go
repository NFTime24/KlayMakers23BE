package logger

import (
	"io"
	"log"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func LogInit(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(warningHandle, "[Debug] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
