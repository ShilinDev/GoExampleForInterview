package els

import (
	"GoExampleForInterview/logger"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"os"
	"time"
)

const (
	ServiceIndex = "consumer"
	ServiceName  = "Go_Consumer"
)

type ElsLogger struct {
	logger *logrus.Logger
}

func (l *ElsLogger) Warn(args ...interface{}) {
	l.logger.Warn(args)
}

func (l *ElsLogger) Error(args ...interface{}) {
	l.logger.Error(args)
}

func (l *ElsLogger) Warning(args ...interface{}) {
	l.logger.Warning(args)
}

func (l *ElsLogger) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l *ElsLogger) Panic(args ...interface{}) {
	l.logger.Panic(args)
}

func (l *ElsLogger) CloseConnection() {
	if l.logger == nil {
		return
	}

	l.logger.Exit(0)
}

func GetLogger() logger.Logger {
	client := &ElsLogger{logger: logrus.New()}

	clientElastic, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_URL")), elastic.SetSniff(false))
	if err != nil {
		panic("Cannot connect to ES")
	}

	// Создание хуков для логирования в эластик
	warnErrLevel, _ := logrus.ParseLevel("warning")
	errorErrLevel, _ := logrus.ParseLevel("error")

	esIndex := fmt.Sprintf("%s-%s", ServiceIndex, time.Now().Format("2006-01-02"))

	warningHook, _ := elogrus.NewAsyncElasticHook(clientElastic, ServiceName, warnErrLevel, esIndex)
	errorHook, _ := elogrus.NewAsyncElasticHook(clientElastic, ServiceName, errorErrLevel, esIndex)

	// навешиваем хуки для записи в эластик ошибок двух уровней
	client.logger.Hooks.Add(warningHook)
	client.logger.Hooks.Add(errorHook)

	return client
}
