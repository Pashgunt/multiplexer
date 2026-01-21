package kafkaconnection

import "time"

const (
	defaultRetryCount = 3
	defaultWorkerCount
)

type Kafka struct {
	timeout      time.Duration
	retryCount   int
	retryTimeout time.Duration
	workerCount  int
}

func (k *Kafka) RetryTimeout() time.Duration {
	return k.retryTimeout
}

func (k *Kafka) Timeout() time.Duration {
	return k.timeout
}

func (k *Kafka) RetryCount() int {
	return k.retryCount
}

func (k *Kafka) WorkerCount() int {
	return k.workerCount
}

func DefaultKafkaConn() Kafka {
	return Kafka{
		timeout:      1 * time.Second,
		retryCount:   defaultRetryCount,
		workerCount:  defaultWorkerCount,
		retryTimeout: 2 * time.Second,
	}
}
