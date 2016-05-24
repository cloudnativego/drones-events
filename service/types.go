package service

import dronescommon "github.com/cloudnativego/drones-common"

const (
	alertsQueueName    = "alerts"
	telemetryQueueName = "telemetry"
	positionsQueueName = "positions"
)

type eventRepository interface {
	UpdateLastTelemetryEvent(telemetryEvent dronescommon.TelemetryUpdatedEvent) (err error)
	UpdateLastAlertEvent(alertEvent dronescommon.AlertSignalledEvent) (err error)
	UpdateLastPositionEvent(positionEvent dronescommon.PositionChangedEvent) (err error)
}
