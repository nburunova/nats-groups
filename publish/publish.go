package main

import (
	"flag"
	"io"
	"os"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	natsURL = "http://localhost:4222"
	subject = "alerts.subject"
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
	msg := "test message"
	logger.Info("try to publish")
	_, err = js.Publish(subject, []byte(msg))
	if err != nil {
		logger.Fatal("failed to publish message",
			zap.String("subject name", subject),
			zap.Any("message", msg),
			zap.Error(err),
		)
	}
	logger.Info("message published", zap.String("subject name", subject), zap.Any("message", msg))

}
