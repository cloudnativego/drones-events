package mongo

import (
	"time"

	dronescommon "github.com/cloudnativego/drones-common"
	"gopkg.in/mgo.v2/bson"
)

func convertAlertEventToRecord(event dronescommon.AlertSignalledEvent, recordID bson.ObjectId) (record *mongoAlertRecord) {
	t := time.Unix(event.ReceivedOn, 0)
	record = &mongoAlertRecord{
		RecordID:    recordID,
		DroneID:     event.DroneID,
		FaultCode:   event.FaultCode,
		Description: event.Description,
		ReceivedOn:  t.Format(time.UnixDate),
	}

	return
}

func convertTelemetryEventToRecord(event dronescommon.TelemetryUpdatedEvent, recordID bson.ObjectId) (record *mongoTelemetryRecord) {
	t := time.Unix(event.ReceivedOn, 0)
	record = &mongoTelemetryRecord{
		RecordID:         recordID,
		DroneID:          event.DroneID,
		RemainingBattery: event.RemainingBattery,
		Uptime:           event.Uptime,
		CoreTemp:         event.CoreTemp,
		ReceivedOn:       t.Format(time.UnixDate),
	}

	return
}

func convertPositionEventToRecord(event dronescommon.PositionChangedEvent, recordID bson.ObjectId) (record *mongoPositionRecord) {
	t := time.Unix(event.ReceivedOn, 0)
	record = &mongoPositionRecord{
		RecordID:        recordID,
		DroneID:         event.DroneID,
		Latitude:        event.Latitude,
		Longitude:       event.Longitude,
		Altitude:        event.Altitude,
		CurrentSpeed:    event.CurrentSpeed,
		HeadingCardinal: event.HeadingCardinal,
		ReceivedOn:      t.Format(time.UnixDate),
	}
	return
}

func convertTelemetryRecordToEvent(record mongoTelemetryRecord) (event dronescommon.TelemetryUpdatedEvent) {
	t, _ := time.Parse(time.UnixDate, record.ReceivedOn)
	event = dronescommon.TelemetryUpdatedEvent{
		DroneID:          record.DroneID,
		RemainingBattery: record.RemainingBattery,
		Uptime:           record.Uptime,
		CoreTemp:         record.CoreTemp,
		ReceivedOn:       t.Unix(),
	}
	return
}

func convertPositionRecordToEvent(record mongoPositionRecord) (event dronescommon.PositionChangedEvent) {
	t, _ := time.Parse(time.UnixDate, record.ReceivedOn)
	event = dronescommon.PositionChangedEvent{
		DroneID:         record.DroneID,
		Altitude:        record.Altitude,
		CurrentSpeed:    record.CurrentSpeed,
		HeadingCardinal: record.HeadingCardinal,
		Latitude:        record.Latitude,
		Longitude:       record.Longitude,
		ReceivedOn:      t.Unix(),
	}
	return
}

func convertAlertRecordToEvent(record mongoAlertRecord) (event dronescommon.AlertSignalledEvent) {
	t, _ := time.Parse(time.UnixDate, record.ReceivedOn)
	event = dronescommon.AlertSignalledEvent{
		DroneID:     record.DroneID,
		Description: record.Description,
		FaultCode:   record.FaultCode,
		ReceivedOn:  t.Unix(),
	}
	return
}
