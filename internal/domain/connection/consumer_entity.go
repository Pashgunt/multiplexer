package kafkaconnection

import "time"

type Consumer struct {
	timeout          time.Duration
	timeoutCommit    time.Duration
	retryCountCommit int
}

func (c Consumer) TimeoutCommit() time.Duration {
	return c.timeoutCommit
}

func (c Consumer) RetryCountCommit() int {
	return c.retryCountCommit
}

func (c Consumer) Timeout() time.Duration {
	return c.timeout
}

func DefaultConsumer() Consumer {
	return Consumer{
		timeout:          3 * time.Second,
		timeoutCommit:    3 * time.Second,
		retryCountCommit: 1,
	}
}
