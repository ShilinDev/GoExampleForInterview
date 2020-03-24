package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TotalMessagesConsumed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "nats_messages_consumed_total",
	Help: "The total number of messages consumed from nats",
})

var TotalParseErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "nats_parse_errors_total",
	Help: "The total number of messages that failed to parse",
})

var TotalPostgresErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "postgresql_insert_errors_total",
	Help: "The total number of requests that failed to insert into postresql",
})

var TotalStatusRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_status_requests_total",
	Help: "The total number of HTTP service status requests",
})

var TotalInvalidRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_invalid_requests_total",
	Help: "The total number of unknown HTTP requests or requests with errors",
})

var ConnectionCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "postgresql_connections",
	Help: "Count of used postgres connections",
})

var SearchQueueCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "queue_search_count",
	Help: "Count of messages in cache search queue",
})

var RoomQueueCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "queue_room_count",
	Help: "Count of messages in room search queue",
})

var ProcessingWorkerTimeVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "process_time_seconds",
		Help: "Amount of time spent processing jobs",
	},
	[]string{"worker_id", "type"},
)

var ProcessingUnmarshalTimeVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "process_unmarshal_time_seconds",
		Help: "Amount of time spent processing for unmarshal jobs",
	},
	[]string{"worker_id", "type"},
)

var ProcessingDBTimeVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "process_db_time_seconds",
		Help: "Amount of time spent processing for unmarshal jobs",
	},
	[]string{"worker_id", "type"},
)

var DBSearchGaugeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "db_search_batch_size",
	Help: "Size in uniq search values in storage",
})

var DBRoomGaugeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "db_room_batch_size",
	Help: "Size in uniq room values in storage",
})

var WorkingWorkersCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "working_workers_count",
	Help: "Count of workers who works right now",
})