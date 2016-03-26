package service

import (
	"fmt"

	common "github.com/cloudnativego/drones-common"
)

func consumeEvents(alertChannel chan common.AlertSignalledEvent, telemetryChannel chan common.TelemetryUpdatedEvent, positionChannel chan common.PositionChangedEvent, repo eventRepository) {
	go func() {
		fmt.Println("Started event consumer goroutine")
		for {
			select {
			case alert := <-alertChannel:
				processAlert(repo, alert)
			case telemetry := <-telemetryChannel:
				processTelemetry(repo, telemetry)
			case position := <-positionChannel:
				processPosition(repo, position)
			}
		}

	}()
}

func processAlert(repo eventRepository, alertEvent common.AlertSignalledEvent) {
	fmt.Printf("Processing alert %+v\n", alertEvent)

	repo.UpdateLastAlertEvent(alertEvent)
	Stats.AlertEventCount++
}

func processTelemetry(repo eventRepository, telemetryEvent common.TelemetryUpdatedEvent) {
	fmt.Printf("Processing telemetry %+v\n", telemetryEvent)
	repo.UpdateLastTelemetryEvent(telemetryEvent)
	Stats.TelemetryEventCount++
}

func processPosition(repo eventRepository, positionEvent common.PositionChangedEvent) {
	fmt.Printf("Processing position %+v\n", positionEvent)
	repo.UpdateLastPositionEvent(positionEvent)
	Stats.PositionEventCount++
}
