package xhandlers

import handlers "learn-go/http/handlers"

type XHandlers struct {
	StudentsHandlers *handlers.StudentsHandler
	LogsHandlers     *handlers.LogsHandler
}
