package nats

import (
	"GoExampleForInterview/app"
	"GoExampleForInterview/config"
	"GoExampleForInterview/consumer"
	"GoExampleForInterview/metrics"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"strconv"
	"time"
)

const (
	SearchSubject = "cache"
	SearchWorker  = "search_worker"
)

type NatsConsumer struct {
	Consumer        *nats.Conn
	searchQueueSize int
}

func (n *NatsConsumer) CloseConnection() {
	if n.Consumer == nil {
		return
	}

	n.Consumer.Close()
}

func (n *NatsConsumer) QueueSubscribe() {
	if _, err := n.Consumer.QueueSubscribe(SearchSubject, SearchWorker, func(m *nats.Msg) {
		if len(app.SearchQueue) == n.searchQueueSize {
			return
		}

		work := app.Job{Message: m.Data}
		app.SearchQueue <- work
		metrics.TotalMessagesConsumed.Inc()
		metrics.SearchQueueCount.Inc()
	}); err != nil {
		config.Instance.Logger.Error(err)
	}
}

func (n *NatsConsumer) GracefulShutdown() {
	if err := n.Consumer.Drain(); err != nil {
		config.Instance.Logger.Error("Couldn't make graceful shutdown")
	}
}

func GetNatsConsumer() consumer.Consumer {
	connectionString := fmt.Sprintf("%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT"))
	opts := []nats.Option{nats.Name("NATS Subscriber")}
	opts = append(opts, nats.UserInfo(os.Getenv("NATS_USER"), os.Getenv("NATS_PASS")))
	opts = setupConnOptions(opts)

	connection, err := nats.Connect(connectionString, opts...)
	if err != nil {
		panic("Couldn't connect to Nats")
	}

	size, err := strconv.Atoi(os.Getenv("SEARCH_QUEUE_SIZE"))
	if err != nil {
		config.Instance.Logger.Warning("Empty QUEUE_SIZE ! Set default value = 100")
		size = 100
	}

	return &NatsConsumer{Consumer: connection, searchQueueSize: size}
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		config.Instance.Logger.Warning("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		config.Instance.Logger.Warning("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		config.Instance.Logger.Error("Exiting: %v", nc.LastError())
	}))
	opts = append(opts, nats.ErrorHandler(natsErrHandler))

	return opts
}

func natsErrHandler(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
	config.Instance.Logger.Error("error: %v\n", natsErr)
	if natsErr == nats.ErrSlowConsumer {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			config.Instance.Logger.Warning("couldn't get pending messages: %v", err)
			return
		}
		config.Instance.Logger.Warning("Falling behind with %d pending messages on subject %q.\n",
			pendingMsgs, sub.Subject)
	}
}
