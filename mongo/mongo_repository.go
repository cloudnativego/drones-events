package mongo

import (
	"errors"
	"time"

	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	dronescommon "github.com/cloudnativego/drones-common"
	"gopkg.in/mgo.v2/bson"
)

/*
type TelemetryUpdatedEvent struct {
	DroneID          string `json:"drone_id"`
	RemainingBattery int    `json:"battery"`
	Uptime           int    `json:"uptime"`
	CoreTemp         int    `json:"core_temp"`
	ReceivedOn       int64  `json:"received_on"`
}

// AlertSignalledEvent is an event indicating an alert condition was reported by a drone
type AlertSignalledEvent struct {
	DroneID     string `json:"drone_id"`
	FaultCode   int    `json:"fault_code"`
	Description string `json:"description"`
	ReceivedOn  int64  `json:"received_on"`
}

// PositionChangedEvent is an event indicating that the position and speed of a drone was reported.
type PositionChangedEvent struct {
	DroneID         string  `json:"drone_id"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	Altitude        float32 `json:"altitude"`
	CurrentSpeed    float32 `json:"current_speed"`
	HeadingCardinal int     `json:"heading_cardinal"`
	ReceivedOn      int64   `json:"received_on"`
}*/

type EventRollupRepository struct {
	PositionsCollection cfmgo.Collection
	TelemetryCollection cfmgo.Collection
	AlertsCollection    cfmgo.Collection
}

type mongoTelemetryRecord struct {
	RecordID         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DroneID          string        `bson:"drone_id",json:"drone_id"`
	RemainingBattery int           `bson:"remaining_battery",json:"remaining_battery"`
	Uptime           int           `bson:"uptime",json:"uptime"`
	CoreTemp         int           `bson:"core_temp",json:"core_temp"`
	ReceivedOn       string        `bson:"received_on",json:"received_on"`
}

type mongoAlertRecord struct {
}

type mongoPositionRecord struct {
}

func NewEventRollupRepository(positions cfmgo.Collection, alerts cfmgo.Collection, telemetry cfmgo.Collection) (repo *EventRollupRepository) {
	repo = &EventRollupRepository{
		PositionsCollection: positions,
		AlertsCollection:    alerts,
		TelemetryCollection: telemetry,
	}
	return
}

func (repo *EventRollupRepository) UpdateLastTelemetryEvent(telemetryEvent dronescommon.TelemetryUpdatedEvent) (err error) {
	repo.TelemetryCollection.Wake()

	var recordID bson.ObjectId
	foundEvent, err := repo.getTelemetryRecord(telemetryEvent.DroneID)
	if err != nil {
		recordID = bson.NewObjectId()
	} else {
		recordID = foundEvent.RecordID
	}

	newRecord := convertTelemetryEventToRecord(telemetryEvent, recordID)
	_, err = repo.TelemetryCollection.UpsertID(recordID, newRecord)

	return
}

func (repo *EventRollupRepository) UpdateLastAlertEvent(alertEvent dronescommon.AlertSignalledEvent) (err error) {
	return
}

func (repo *EventRollupRepository) UpdateLastPositionEvent(positionEvent dronescommon.PositionChangedEvent) (err error) {
	return
}

func (repo *EventRollupRepository) getTelemetryRecord(droneID string) (record mongoTelemetryRecord, err error) {
	var records []mongoTelemetryRecord
	query := bson.M{"drone_id": droneID}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.TelemetryCollection.Find(params, &records)
	if count == 0 {
		err = errors.New("Telemetry record not found.")
	}
	if err == nil {
		record = records[0]
	}
	return
}

func convertTelemetryEventToRecord(event dronescommon.TelemetryUpdatedEvent, recordID bson.ObjectId) (record *mongoTelemetryRecord) {
	time := time.Unix(event.ReceivedOn, 0)
	record = &mongoTelemetryRecord{
		RecordID:         recordID,
		DroneID:          event.DroneID,
		RemainingBattery: event.RemainingBattery,
		Uptime:           event.Uptime,
		CoreTemp:         event.CoreTemp,
		ReceivedOn:       time.Format("2006-01-02 15:04:05"),
	}

	return
}
