package fslnpkg

import (
	"bitbucket.org/fstlnrd/fstln-package/pkg/mq/rabbit/consumer"
	"testing"
)

func Test1(t *testing.T) {
	consumer.CreateRabbitConsumer()
}
