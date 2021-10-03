package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type consumerCfg struct {
	name          string
	subject       string
	group         string
	durable       string
	streamName    string
	streamSubject string
}

var (
	natsURL   = "http://localhost:4222"
	consumer1 = consumerCfg{
		name:          "consumer1",
		subject:       "alerts.subject1",
		group:         "group1",
		durable:       "durable1",
		streamName:    "alerts",
		streamSubject: "alerts.*",
	}
	consumer2 = consumerCfg{
		name:          "consumer2",
		subject:       "alerts.subject2",
		group:         "group2",
		durable:       "durable2",
		streamName:    "alerts",
		streamSubject: "alerts.*",
	}
	deliverPolicyString = "last"
)

type message struct {
	ID        int
	Timestamp time.Time
}

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

	deliverPolicy, deliverPolicyOpt := deliverPolicyFromString(deliverPolicyString)

	done := make(chan struct{})
	catched := make(chan struct{})
	defer close(done)

	subs := make([]*nats.Subscription, 0)

	for _, c := range []consumerCfg{consumer2, consumer1} {
		consumer := c
		stream, err := js.StreamInfo(consumer.streamName)
		if err != nil {
			logger.Info("failed to get stream", zap.String("stream name", consumer.streamName), zap.Error(err))
		}
		if stream == nil {
			stream, err = js.AddStream(&nats.StreamConfig{
				Name:     consumer.streamName,
				Subjects: []string{consumer.streamSubject},
			})
			if err != nil {
				logger.Fatal("failed to add stream",
					zap.String("stream name", consumer.streamName),
					zap.String("stream subject", consumer.streamSubject),
					zap.Error(err),
				)
			} else {
				logger.Info("stream added", zap.Any("stream", stream))
			}
		}

		ci, err := js.ConsumerInfo(consumer.streamName, consumer.name)
		if err != nil {
			logger.Info("failed to get consumer",
				zap.String("stream name", consumer.streamName),
				zap.String("consumer name", consumer.name),
				zap.Error(err),
			)
		}
		if ci == nil {
			ci, err = js.AddConsumer(consumer.streamName, &nats.ConsumerConfig{
				AckPolicy:      nats.AckExplicitPolicy,
				Durable:        consumer.durable,
				AckWait:        5 * time.Second,
				DeliverSubject: consumer.group,
				DeliverGroup:   consumer.group,
				DeliverPolicy:  deliverPolicy,
			})
			if err != nil {
				logger.Fatal("failed to add consumer",
					zap.String("stream name", consumer.streamName),
					zap.String("consumer name", consumer.name),
					zap.Error(err),
				)
			} else {
				logger.Info("consumer added", zap.Any("consumer", ci), zap.Error(err))
			}
		}

		sub, err := js.QueueSubscribe(consumer.subject, consumer.group, func(m *nats.Msg) {
			var msg message
			if err := json.Unmarshal(m.Data, &msg); err != nil {
				logger.Error("failed to unmarshal message", zap.ByteString("raw", m.Data), zap.Error(err))
				return
			}
			if err := m.Ack(); err != nil {
				logger.Error("failed to send ack", zap.Any("message", msg), zap.Error(err))
				return
			}
			catched <- struct{}{}
			logger.Info(fmt.Sprintf("subject %v, group %v: msg received", consumer.subject, consumer.group), zap.Any("message", msg))
		},
			nats.ManualAck(),
			nats.Durable(ci.Config.Durable),
			nats.AckWait(ci.Config.AckWait),
			nats.DeliverSubject(ci.Config.DeliverSubject),
			nats.Bind(consumer.streamName, ci.Name),
			deliverPolicyOpt,
		)

		if err != nil {
			logger.Fatal("failed to subscribe",
				zap.String("subject name", consumer.subject),
				zap.String("queue name", consumer.group),
				zap.Error(err),
			)
		} else {
			logger.Info("subscription added", zap.Any("sub", sub), zap.Error(err))
		}

		subs = append(subs, sub)
	}

	for i, subject := range []string{consumer1.subject, consumer2.subject} {
		msg := &message{
			ID:        int(i),
			Timestamp: time.Now(),
		}
		b, err := json.Marshal(msg)
		if err != nil {
			logger.Fatal("failed to marshal message", zap.Any("message", msg), zap.Error(err))
		}
		_, err = js.Publish(subject, b)
		if err != nil {
			logger.Fatal("failed to publish message",
				zap.String("subject name", subject),
				zap.Any("message", msg),
				zap.Error(err),
			)
		}
		logger.Info("message published", zap.String("subject name", subject), zap.Any("message", msg))
	}

	timeout := 5 * time.Second
	ticker := time.NewTicker(timeout)

	for {
		select {
		case <-ticker.C:
			logger.Fatal("===== TIMED OUT =====")
		case <-catched:
			ticker.Reset(timeout)
		}
	}
	for _, sub := range subs {
		sub.Drain()
	}
	nc.Drain()
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
