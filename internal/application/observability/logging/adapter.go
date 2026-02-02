package logging

import (
	"sync"
	"transport/pkg/utils/backoff"
)

type AdapterInterface interface {
	GetLogger(loggerType backoff.LoggerType) LoggerInterface
}

type Adapter struct {
	loggers map[backoff.LoggerType]LoggerInterface
	once    sync.Once
	mutex   sync.RWMutex
	factory LoggerFactory
}

func (adapter *Adapter) GetLogger(loggerType backoff.LoggerType) LoggerInterface {
	adapter.mutex.RLock()
	defer adapter.mutex.RUnlock()
	logger, ok := adapter.loggers[loggerType]

	if !ok {
		adapter.mutex.RUnlock()
		adapter.mutex.Lock()

		logger = adapter.factory.CreateLogger(loggerType)
		adapter.loggers[loggerType] = logger

		adapter.mutex.Unlock()
		adapter.mutex.RLock()
	}

	return logger
}

func (adapter *Adapter) Init(loggerTypes []backoff.LoggerType) {
	adapter.once.Do(func() {
		adapter.mutex.Lock()
		defer adapter.mutex.Unlock()

		for _, loggerType := range loggerTypes {
			if _, ok := adapter.loggers[loggerType]; ok {
				continue
			}

			adapter.loggers[loggerType] = adapter.factory.CreateLogger(loggerType)
		}
	})
}

func NewAdapter() AdapterInterface {
	return &Adapter{
		factory: &defaultLoggerFactory{},
	}
}
