package logger

type Logger interface {
	Warn(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Info(args ...interface{})
	Panic(args ...interface{})
	CloseConnection()
}
