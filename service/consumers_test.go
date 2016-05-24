package service

import (
	"fmt"
	"math"
	"testing"
	"time"

	dronescommon "github.com/cloudnativego/drones-common"
)

func TestConsumer(t *testing.T) {
	alertChannel := make(chan dronescommon.AlertSignalledEvent)
	telemetryChannel := make(chan dronescommon.TelemetryUpdatedEvent)
	positionChannel := make(chan dronescommon.PositionChangedEvent)

	repo := &fakeRepository{}
	consumeEvents(alertChannel, telemetryChannel, positionChannel, repo)
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Sending messages...")
		alertChannel <- dronescommon.AlertSignalledEvent{DroneID: "abc", FaultCode: 1, Description: "Faily fail"}
		alertChannel <- dronescommon.AlertSignalledEvent{DroneID: "abc2", FaultCode: 2, Description: "More fails"}
		telemetryChannel <- dronescommon.TelemetryUpdatedEvent{DroneID: "abc", RemainingBattery: 12, Uptime: 10, CoreTemp: 10}
		positionChannel <- dronescommon.PositionChangedEvent{DroneID: "abc", Latitude: 32.45, Longitude: 75.212, Altitude: 105.3, CurrentSpeed: 17.32}
	}()

	// HACK
	time.Sleep(time.Second * 4)

	if repo.lastAlert.DroneID != "abc2" {
		t.Errorf("Last alert was not 'abc2', got '%s'", repo.lastAlert.DroneID)
	}

	if repo.lastTelemetry.CoreTemp != 10 {
		t.Errorf("Expected last telemetry core temp of 10, got %d", repo.lastTelemetry.CoreTemp)
	}

	if math.Floor(float64(repo.lastPosition.Latitude)) != 32 {
		t.Errorf("Expected last position latitude of 32.45, got %f", repo.lastPosition.Latitude)
	}
}
