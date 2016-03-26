package service

import dronescommon "github.com/cloudnativego/drones-common"

type fakeRepository struct {
	lastAlert     dronescommon.AlertSignalledEvent
	lastPosition  dronescommon.PositionChangedEvent
	lastTelemetry dronescommon.TelemetryUpdatedEvent
}

func NewFakeRepository() *fakeRepository {
	return &fakeRepository{}
}

func (repo *fakeRepository) UpdateLastTelemetryEvent(telemetryEvent dronescommon.TelemetryUpdatedEvent) (err error) {

	repo.lastTelemetry = telemetryEvent
	return
}

func (repo *fakeRepository) UpdateLastAlertEvent(alertEvent dronescommon.AlertSignalledEvent) (err error) {

	repo.lastAlert = alertEvent
	return
}

func (repo *fakeRepository) UpdateLastPositionEvent(positionEvent dronescommon.PositionChangedEvent) (err error) {

	repo.lastPosition = positionEvent
	return
}
