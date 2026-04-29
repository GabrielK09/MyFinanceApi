package loggerHelper

import (
	"log"
	"os"
)

var GeneralLogger *log.Logger
var InfoLogger *log.Logger
var ErrorLogger *log.Logger
var SuccessLogger *log.Logger

func Logger() {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "dev" {
		GeneralLogger = newFileLogger("General Logger:\t")
		ErrorLogger = newFileLogger("Error Logger:\t")
		InfoLogger = newFileLogger("Info Logger:\t")
		SuccessLogger = newFileLogger("Success Logger:\t")
		return
	}

	GeneralLogger = log.New(os.Stdout, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "Info Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	SuccessLogger = log.New(os.Stdout, "Success Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)

}

func newFileLogger(prefix string) *log.Logger {
	logDir := "log"

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatal("Erro ao criar o diretório de log: ", err)
	}

	file, err := os.OpenFile("log/general-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Erro ao abrir o log:", err)
	}

	return log.New(file, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}
