package kafkaconnection

import "time"

const (
	defaultRetryCount = 3
	defaultWorkerCount
)

type Kafka struct {
	timeout     time.Duration
	retryCount  int
	workerCount int
}

func (k *Kafka) Timeout() time.Duration {
	return k.timeout
}

func (k *Kafka) SetTimeout(timeout time.Duration) {
	k.timeout = timeout
}

func (k *Kafka) RetryCount() int {
	return k.retryCount
}

func (k *Kafka) SetRetryCount(retryCount int) {
	k.retryCount = retryCount
}

func (k *Kafka) WorkerCount() int {
	return k.workerCount
}

func (k *Kafka) SetWorkerCount(workerCount int) {
	k.workerCount = workerCount
}

func DefaultKafkaConn() Kafka {
	return Kafka{
		timeout:     1 * time.Second,
		retryCount:  defaultRetryCount,
		workerCount: defaultWorkerCount,
	}
}
