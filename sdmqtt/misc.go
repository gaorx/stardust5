package sdmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Loggers struct {
	Debug    mqtt.Logger
	Error    mqtt.Logger
	Warn     mqtt.Logger
	Critical mqtt.Logger
}

func SetLogger(loggers Loggers) {
	if loggers.Debug != nil {
		mqtt.DEBUG = loggers.Debug
	}
	if loggers.Error != nil {
		mqtt.ERROR = loggers.Error
	}
	if loggers.Warn != nil {
		mqtt.WARN = loggers.Warn
	}
	if loggers.Critical != nil {
		mqtt.CRITICAL = loggers.Critical
	}
}
