package consts

type contextKey string

const (
	RequestIDKey  contextKey = "request_id"
	RequestLogger contextKey = "request_logger"
)
