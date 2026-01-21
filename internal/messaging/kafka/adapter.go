package kafka

import (
	"context"
	"sync"
	"time"
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/infrastructure/config/types"
)

type Adapter struct {
	connections []*Connection
	configs     []Config
	mutex       sync.RWMutex
}

func NewAdapter(appConfig types.Config) *Adapter {
	adapter := &Adapter{
		configs: convert(appConfig),
		mutex:   sync.RWMutex{},
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

	for i := 0; i < kafka.WorkerCount(); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			for config := range configsChan {
				connection, err := adapter.doConnect(config, kafka.Timeout(), kafka.RetryCount())

				if err != nil {
					continue
				}

				resultsChan <- connection
			}
		}(i)
	}

	adapter.fillConfigs(configsChan)

	close(configsChan)

	go func() {
		wg.Wait()
		close(resultsChan)
	}()
}

func (adapter *Adapter) doConnect(config Config, timeout time.Duration, retryCount int) (*Connection, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connection, err := NewConnection(connCtx, config)

	if err != nil && retryCount > 0 {
		return adapter.doConnect(config, timeout, retryCount-1)
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
