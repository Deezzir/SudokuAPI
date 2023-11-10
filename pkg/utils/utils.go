package utils

import (
	"log"
	"os"
)

var ErrorLog *log.Logger
var InfoLog *log.Logger
var WarningLog *log.Logger

func init() {
	log_flags := log.LstdFlags | log.Lshortfile
	InfoLog = log.New(os.Stdout, "[INFO]:", log_flags)
	ErrorLog = log.New(os.Stderr, "[ERROR]:", log_flags)
	WarningLog = log.New(os.Stdout, "[WARNING]:", log_flags)
}
