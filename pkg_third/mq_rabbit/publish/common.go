package send

import (
	"context"
	"log"
	"os"
	"pkg_third/mq_rabbit/rbmq"
	"strings"
)

var ctx = context.Background()
var conn = rbmq.NewConn("amqp://root:123456@localhost:5672/")

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
