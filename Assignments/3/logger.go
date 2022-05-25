package main

import (
	"crypto/tls"
	"fmt"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

func NewLogrus(url, port string) *logrus.Logger {
	logger := logrus.New()
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", url, port), &tls.Config{RootCAs: nil})
	if err != nil {
		logger.Fatalf("failed to connect to logstash: %v", err)
	}
	logger.Hooks.Add(logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "law-assignment-3"})))

	return logger
}
