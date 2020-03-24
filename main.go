// Пакет для работы с брокером сообщений Nats
package main

import (
	"GoExampleForInterview/config"
	"GoExampleForInterview/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	config.Init()
}

func main() {
	// Собираем список интересующих нас сигналов
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Запуск серверов
	startServers()
	// Инициализация метрик Prometheus
	startMetrics()

	// Подписка на канал натса
	config.Instance.Consumer.QueueSubscribe()

	for {
		select {
		// кейс в случае выключения внешним сигналом
		case <-signals:
			config.Instance.Logger.Warning("Interrupted")
			config.Instance.Consumer.GracefulShutdown()
			os.Exit(0)
		// кейс в случае ошибки внутри канала
		case err := <-app.ErrorsChannel:
			config.Instance.Logger.Error(err)
			metrics.TotalParseErrors.Inc()
		}
	}
}

//инициализация метрик для прометеуса
func startMetrics() {
	prometheus.MustRegister(metrics.SearchQueueCount)
	prometheus.MustRegister(metrics.RoomQueueCount)
	prometheus.MustRegister(metrics.ProcessingWorkerTimeVec)
	prometheus.MustRegister(metrics.ProcessingDBTimeVec)
	prometheus.MustRegister(metrics.ProcessingUnmarshalTimeVec)
	prometheus.MustRegister(metrics.ConnectionCount)
	prometheus.MustRegister(metrics.DBRoomGaugeMetric)
	prometheus.MustRegister(metrics.DBSearchGaugeMetric)
	prometheus.MustRegister(metrics.WorkingWorkersCount)
}

// запуск серверов c роутингом
func startServers() {
	router := mux.NewRouter()
	router.HandleFunc("/", statusHandler).Methods("GET")

	goPprof := router.PathPrefix("/debug/pprof").Subrouter()
	goPprof.HandleFunc("/", pprof.Index)
	goPprof.HandleFunc("/cmdline", pprof.Cmdline)
	goPprof.HandleFunc("/symbol", pprof.Symbol)
	goPprof.HandleFunc("/trace", pprof.Trace)

	profile := goPprof.PathPrefix("/profile").Subrouter()
	profile.HandleFunc("", pprof.Profile)
	profile.Handle("/goroutine", pprof.Handler("goroutine"))
	profile.Handle("/threadcreate", pprof.Handler("threadcreate"))
	profile.Handle("/heap", pprof.Handler("heap"))
	profile.Handle("/block", pprof.Handler("block"))
	profile.Handle("/mutex", pprof.Handler("mutex"))

	router.NotFoundHandler = app.NotFoundHandler{}
	go http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
	promRouter := mux.NewRouter()
	promRouter.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+os.Getenv("PROMETHEUS_SERVER_PORT"), promRouter)
}

func statusHandler(writer http.ResponseWriter, _ *http.Request) {
	_, _ = writer.Write([]byte(`{"status" : "OK"}`))
	metrics.TotalStatusRequests.Inc()
}
