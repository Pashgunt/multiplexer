package logging

import (
	"sync"
	"transport/pkg/utils/backoff"
)

type AdapterInterface interface {
	GetLogger(loggerType backoff.LoggerType) LoggerInterface
	Init(loggerTypes []backoff.LoggerType)
}

type Adapter struct {
	loggers map[backoff.LoggerType]LoggerInterface
	once    sync.Once
	mutex   sync.RWMutex
	factory LoggerFactory
	levels  map[backoff.LoggerType]backoff.LoggerLevel
}

func (adapter *Adapter) GetLogger(loggerType backoff.LoggerType) LoggerInterface {
	adapter.mutex.RLock()
	defer adapter.mutex.RUnlock()
	logger, ok := adapter.loggers[loggerType]

	if !ok {
		adapter.mutex.RUnlock()
		adapter.mutex.Lock()

		logger = adapter.factory.CreateLogger(loggerType, adapter.levels[loggerType])
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

			adapter.loggers[loggerType] = adapter.factory.CreateLogger(loggerType, adapter.levels[loggerType])
		}
	})
}

func NewAdapter(levels map[backoff.LoggerType]backoff.LoggerLevel) AdapterInterface {
	return &Adapter{
		factory: NewDefaultLoggerFactory(),
		loggers: map[backoff.LoggerType]LoggerInterface{},
		levels:  levels,
	}
}
