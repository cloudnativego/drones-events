package mongo

import (
	"errors"

	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	dronescommon "github.com/cloudnativego/drones-common"
	"gopkg.in/mgo.v2/bson"
)

// NewEventRollupRepository creates a new mongoDB event rollup repository with the supplied collections.
func NewEventRollupRepository(positions cfmgo.Collection, alerts cfmgo.Collection, telemetry cfmgo.Collection) (repo *EventRollupRepository) {
	repo = &EventRollupRepository{
		PositionsCollection: positions,
		AlertsCollection:    alerts,
		TelemetryCollection: telemetry,
	}
	return
}

// UpdateLastTelemetryEvent updates the most recent telemetry event, or creates a new one if one has never been received.
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

// UpdateLastAlertEvent updates the most recent alert event, or creates a new one if one has never been received.
func (repo *EventRollupRepository) UpdateLastAlertEvent(alertEvent dronescommon.AlertSignalledEvent) (err error) {
	repo.AlertsCollection.Wake()

	var recordID bson.ObjectId
	foundEvent, err := repo.getAlertRecord(alertEvent.DroneID)
	if err != nil {
		recordID = bson.NewObjectId()
	} else {
		recordID = foundEvent.RecordID
	}

	newRecord := convertAlertEventToRecord(alertEvent, recordID)
	_, err = repo.AlertsCollection.UpsertID(recordID, newRecord)
	return
}

// UpdateLastPositionEvent updates the last position event, or creates a new one if one has never been received.
func (repo *EventRollupRepository) UpdateLastPositionEvent(positionEvent dronescommon.PositionChangedEvent) (err error) {
	repo.PositionsCollection.Wake()

	var recordID bson.ObjectId
	foundEvent, err := repo.getPositionRecord(positionEvent.DroneID)
	if err != nil {
		recordID = bson.NewObjectId()
	} else {
		recordID = foundEvent.RecordID
	}

	newRecord := convertPositionEventToRecord(positionEvent, recordID)
	_, err = repo.PositionsCollection.UpsertID(recordID, newRecord)
	return
}

// GetTelemetryEvent retrieves the most recent telemetry event for a given drone.
func (repo *EventRollupRepository) GetTelemetryEvent(droneID string) (event dronescommon.TelemetryUpdatedEvent, err error) {
	record, err := repo.getTelemetryRecord(droneID)

	if err == nil {
		event = convertTelemetryRecordToEvent(record)
	}

	return
}

// GetAlertEvent retrieves the most recent alert event for a given drone.
func (repo *EventRollupRepository) GetAlertEvent(droneID string) (event dronescommon.AlertSignalledEvent, err error) {
	record, err := repo.getAlertRecord(droneID)
	if err == nil {
		event = convertAlertRecordToEvent(record)
	}
	return
}

// GetPositionEvent retrieves the most recent position event for a given drone.
func (repo *EventRollupRepository) GetPositionEvent(droneID string) (event dronescommon.PositionChangedEvent, err error) {
	record, err := repo.getPositionRecord(droneID)
	if err == nil {
		event = convertPositionRecordToEvent(record)
	}
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

func (repo *EventRollupRepository) getAlertRecord(droneID string) (record mongoAlertRecord, err error) {
	var records []mongoAlertRecord
	query := bson.M{"drone_id": droneID}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.AlertsCollection.Find(params, &records)
	if count == 0 {
		err = errors.New("Alert record not found.")
	}
	if err == nil {
		record = records[0]
	}
	return
}

func (repo *EventRollupRepository) getPositionRecord(droneID string) (record mongoPositionRecord, err error) {
	var records []mongoPositionRecord
	query := bson.M{"drone_id": droneID}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.PositionsCollection.Find(params, &records)
	if count == 0 {
		err = errors.New("Position record not found.")
	}
	if err == nil {
		record = records[0]
	}
	return
}
