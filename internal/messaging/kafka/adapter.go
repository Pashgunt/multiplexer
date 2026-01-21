package kafka

import (
	"context"
	"sync"
	"time"
	"transport/internal/application/observability/logging"
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/infrastructure/config/types"
)

type Adapter struct {
	connections []*Connection
	configs     []Config
	mutex       sync.RWMutex
	logger      logging.KafkaConnectionLogger
}

func NewAdapter(appConfig types.Config, logger logging.KafkaConnectionLogger) *Adapter {
	adapter := &Adapter{
		configs: convert(appConfig),
		mutex:   sync.RWMutex{},
		logger:  logger,
	}
	adapter.connections = make([]*Connection, 0, len(adapter.configs))

	return adapter
}

func (adapter *Adapter) Connections() []*Connection {
	adapter.mutex.RLock()
	defer adapter.mutex.RUnlock()

	copyConnections := make([]*Connection, len(adapter.configs))
	copy(copyConnections, adapter.connections)

	return copyConnections
}

func (adapter *Adapter) ConnectAll(kafka kafkaconnection.Kafka) {
	configsChan := make(chan Config, len(adapter.configs))
	resultsChan := make(chan *Connection, len(adapter.configs))

	var wg sync.WaitGroup

	go adapter.fillConnections(resultsChan)

	for workerIndex := 0; workerIndex < kafka.WorkerCount(); workerIndex++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for config := range configsChan {
				connection, err := adapter.doConnect(config, kafka.Timeout(), kafka.RetryTimeout(), kafka.RetryCount())

				if err != nil {
					adapter.logger.Info(logging.KafkaConnectionLogEntity{
						Message: "Failed to connect to Kafka. " + err.Error(),
						Broker:  config.Broker,
					})
					continue
				}

				adapter.logger.Info(logging.KafkaConnectionLogEntity{
					Message: "Successfully connected to Kafka.",
					Broker:  config.Broker,
				})
				resultsChan <- connection
			}
		}()
	}

	adapter.fillConfigs(configsChan)

	close(configsChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()
}

func (adapter *Adapter) doConnect(
	config Config,
	timeout time.Duration,
	retryTimeout time.Duration,
	retryCount int,
) (*Connection, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connection, err := NewConnection(connCtx, config)

	if err != nil && retryCount > 0 {
		adapter.logger.Info(logging.KafkaConnectionLogEntity{
			Message: "Failed to connect to Kafka. Retrying... " + err.Error(),
			Broker:  config.Broker,
		})

		select {
		case <-time.After(retryTimeout):
			return adapter.doConnect(config, timeout, retryTimeout, retryCount-1)
		}
	}

	return connection, err
}

func (adapter *Adapter) fillConfigs(configsChan chan<- Config) {
	for _, config := range adapter.configs {
		configsChan <- config
	}
}

func (adapter *Adapter) fillConnections(resultsChan <-chan *Connection) {
	for connection := range resultsChan {
		adapter.mutex.Lock()
		adapter.connections = append(adapter.connections, connection)
		adapter.mutex.Unlock()
	}
}
