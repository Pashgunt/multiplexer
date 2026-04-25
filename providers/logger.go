package providers

import (
	"transport/internal/infrastructure/config"
	"transport/pkg/logging"
	"transport/pkg/utils/backoff"
)

func Logger(env config.EnvironmentInterface) logging.AdapterInterface {
	logger := logging.NewAdapter(map[backoff.LoggerType]backoff.LoggerLevel{
		backoff.KafkaLogger: backoff.LoggerLevel(env.Get(backoff.EnvKafkaDebugLevelKey)),
		backoff.AppLogger:   backoff.LoggerLevel(env.Get(backoff.EnvAppDebugLevelKey)),
		backoff.APILogger:   backoff.LoggerLevel(env.Get(backoff.EnvAPIDebugLevelKey)),
	})
	logger.Init([]backoff.LoggerType{backoff.KafkaLogger, backoff.AppLogger, backoff.APILogger})

	return logger
}
