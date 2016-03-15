package service

import dronescommon "github.com/cloudnativego/drones-common"

type fakeRepository struct {
	lastAlert     dronescommon.AlertSignalledEvent
	lastPosition  dronescommon.PositionChangedEvent
	lastTelemetry dronescommon.TelemetryUpdatedEvent
}

func (repo *fakeRepository) updateLastTelemetryEvent(telemetryEvent dronescommon.TelemetryUpdatedEvent) (err error) {

	repo.lastTelemetry = telemetryEvent
	return
}

func (repo *fakeRepository) updateLastAlertEvent(alertEvent dronescommon.AlertSignalledEvent) (err error) {

	repo.lastAlert = alertEvent
	return
}

func (repo *fakeRepository) updateLastPositionEvent(positionEvent dronescommon.PositionChangedEvent) (err error) {

	repo.lastPosition = positionEvent
	return
}
