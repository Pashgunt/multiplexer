package kafka

import (
	"context"
	"sync"
	"time"
	"transport/internal/application/observability/logging"
	kafkaconnection "transport/internal/domain/connection"
	appconfig "transport/internal/infrastructure/config/app"
	"transport/pkg/utils/backoff"
)

type AdapterInterface interface {
	AdapterGetterInterface
	ConnectAll(kafka kafkaconnection.Kafka)
	CloseAll(ctx context.Context)
}

type AdapterGetterInterface interface {
	Connections() []ConnectionInterface
}

// Adapter todo add flag success connections
type Adapter struct {
	connections []ConnectionInterface
	configs     []Config
	mutex       sync.RWMutex
	logger      logging.LoggerInterface
}

func NewAdapter(config appconfig.Config) AdapterInterface {
	adapter := &Adapter{
		configs: convert(config.Config),
		mutex:   sync.RWMutex{},
		logger:  config.Logger.GetLogger(backoff.KafkaLogger),
	}
	adapter.connections = make([]ConnectionInterface, 0, len(adapter.configs))

	return adapter
}

func (adapter *Adapter) Connections() []ConnectionInterface {
	adapter.mutex.RLock()
	defer adapter.mutex.RUnlock()

	copyConnections := make([]ConnectionInterface, len(adapter.configs))
	copy(copyConnections, adapter.connections)

	return copyConnections
}

func (adapter *Adapter) ConnectAll(kafka kafkaconnection.Kafka) {
	configsChan := make(chan Config, len(adapter.configs))
	resultsChan := make(chan ConnectionInterface, len(adapter.configs))

	var wg sync.WaitGroup

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
	wg.Wait()
	close(resultsChan)
	adapter.fillConnections(resultsChan)
}

func (adapter *Adapter) CloseAll(ctx context.Context) {
	if len(adapter.connections) == 0 {
		return
	}

	errorCloseConnections := make([]ConnectionInterface, 0, len(adapter.connections))

	for index, connection := range adapter.connections {
		select {
		case <-ctx.Done():
			adapter.logger.Info(logging.KafkaConnectionLogEntity{
				Message: "Context canceled during connection closing",
			})
			errorCloseConnections = append(errorCloseConnections, adapter.connections[index:]...)
			adapter.connections = errorCloseConnections
			return
		default:
		}

		if err := connection.Close(); err != nil {
			adapter.logger.Info(logging.KafkaConnectionLogEntity{
				Message: "Error closing Kafka connection: " + err.Error(),
			})
			errorCloseConnections = append(errorCloseConnections, connection)
		}

		adapter.logger.Info(logging.KafkaConnectionLogEntity{
			Message: "Successfully closed Kafka connections.",
		})
	}

	adapter.connections = errorCloseConnections

	if len(adapter.connections) > 0 {
		adapter.CloseAll(ctx)
	}
}

func (adapter *Adapter) doConnect(
	config Config,
	timeout time.Duration,
	retryTimeout time.Duration,
	retryCount int,
) (ConnectionInterface, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connection, err := NewConnection(connCtx, config, adapter.logger)

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

func (adapter *Adapter) fillConnections(resultsChan <-chan ConnectionInterface) {
	for connection := range resultsChan {
		adapter.mutex.Lock()
		adapter.connections = append(adapter.connections, connection)
		adapter.mutex.Unlock()
	}
}
