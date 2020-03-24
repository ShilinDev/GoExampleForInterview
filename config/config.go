package config

import (
	"GoExampleForInterview/consumer"
	"GoExampleForInterview/consumer/nats"
	"GoExampleForInterview/db"
	"GoExampleForInterview/db/postgres"
	"GoExampleForInterview/logger"
	"GoExampleForInterview/logger/els"
	"github.com/joho/godotenv"
)

type InstanceConfig struct {
	Logger   logger.Logger
	Consumer consumer.Consumer
	DB       db.Connection
}

var Instance = InstanceConfig{}

func Init() {
	defer func() {
		if err := recover(); err != nil {
			Instance.Logger.CloseConnection()
			Instance.Consumer.CloseConnection()
			Instance.DB.CloseConnection()
		}
	}()
	// Подключение файла окружения
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Получаем реализацию интерфейса логгера
	Instance.Logger = els.GetLogger()
	// Получаем реализацию интерфейса брокера сообщений
	Instance.Consumer = nats.GetNatsConsumer()
	// Получаем реализацию интерфейса пулла базы данных
	Instance.DB = postgres.GetPostgresPullConnection()
}
