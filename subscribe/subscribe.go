package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	natsURL             = "http://localhost:4222"
	subject             = "alerts.subject"
	queue               = "queue"
	streamName          = "alerts"
	streamSubject       = "alerts.*"
	deliverPolicyString = "last"

	differentDurableNames = false
)

func main() {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(io.Writer(os.Stdout)),
			zap.NewAtomicLevelAt(zapcore.DebugLevel),
		),
		zap.Development(),
		zap.ErrorOutput(zapcore.AddSync(io.Writer(os.Stderr))),
	).With(zap.Int("pid", os.Getpid()))

	flag.StringVar(&natsURL, "url", natsURL, "NATS url")
	flag.BoolVar(&differentDurableNames, "diff", false, "Different durable names for subscriptions")
	flag.Parse()

	nc, err := nats.Connect(natsURL)
	if err != nil {
		logger.Fatal("failed to connect to NATS", zap.Error(err))
	}
	defer nc.Close()
	js, err := nc.JetStream()
	if err != nil {
		logger.Fatal("failed to get JetStream", zap.Error(err))
	}

	_, deliverPolicyOpt := deliverPolicyFromString(deliverPolicyString)

	subs := make([]*nats.Subscription, 0)

	for i := 1; i <= 2; i++ {
		stream, err := js.StreamInfo(streamName)
		if err != nil {
			logger.Info("failed to get stream", zap.String("stream name", streamName), zap.Error(err))
		}
		if stream == nil {
			stream, err = js.AddStream(&nats.StreamConfig{
				Name:     streamName,
				Subjects: []string{streamSubject},
			})
			if err != nil {
				logger.Fatal("failed to add stream",
					zap.String("stream name", streamName),
					zap.String("stream subject", streamSubject),
					zap.Error(err),
				)
			} else {
				logger.Info("stream added", zap.Any("stream", stream))
			}
		}
		durable := "subscription"
		if differentDurableNames {
			durable = fmt.Sprintf("subscription_%v", i)
		}
		sub, err := js.QueueSubscribe(subject, queue, func(m *nats.Msg) {
			if err := m.Ack(); err != nil {
				logger.Error("failed to send ack", zap.Any("message", m), zap.Error(err))
				return
			}
			logger.Info(fmt.Sprintf("subject %v, group %v: msg received", subject, queue), zap.Any("message", m))
		},
			nats.ManualAck(),
			nats.Durable(durable),
			nats.AckWait(5*time.Second),
			deliverPolicyOpt,
		)

		if err != nil {
			logger.Fatal("failed to subscribe",
				zap.String("subject name", subject),
				zap.String("queue name", queue),
				zap.Error(err),
			)
		} else {
			logger.Info("subscription added", zap.Any("durable", durable), zap.Any("sub", sub), zap.Error(err))
		}

		subs = append(subs, sub)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	for _, sub := range subs {
		sub.Drain()
	}
	nc.Drain()
	log.Fatalf("Exiting")
}

func deliverPolicyFromString(s string) (nats.DeliverPolicy, nats.SubOpt) {
	switch strings.ToLower(s) {
	case "all":
		return nats.DeliverAllPolicy, nats.DeliverAll()
	case "new":
		return nats.DeliverNewPolicy, nats.DeliverNew()
	case "last":
		fallthrough
	default:
		return nats.DeliverLastPolicy, nats.DeliverLast()
	}
}
