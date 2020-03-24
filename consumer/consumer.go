package consumer

type Consumer interface {
	QueueSubscribe()
	CloseConnection()
	GracefulShutdown()
}
