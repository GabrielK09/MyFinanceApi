package handler

import loggerHelper "my_finance/internal/logger"

func HandlePanic() {
	a := recover()

	if a != nil {
		loggerHelper.ErrorLogger.Println("PANIC: ", a)
	}
}
