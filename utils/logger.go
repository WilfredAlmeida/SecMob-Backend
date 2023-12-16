package utils

import (
	"fmt"
	"log"
)

type AggregatedLogger struct {
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

var MLogger = AggregatedLogger{}

func (l *AggregatedLogger) InfoLog(v ...interface{}) {
	_ = l.InfoLogger.Output(2, fmt.Sprintln(v...))
}
func (l *AggregatedLogger) WarnLog(v ...interface{}) {
	_ = l.WarnLogger.Output(2, fmt.Sprintln(v...))
}
func (l *AggregatedLogger) ErrorLog(v ...interface{}) {
	_ = l.ErrorLogger.Output(2, fmt.Sprintln(v...))
}
