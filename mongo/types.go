package mongo

import (
	"github.com/cloudnativego/cfmgo"
	"gopkg.in/mgo.v2/bson"
)

// EventRollupRepository is the anchor struct for mongoDB repository implementation methods.
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
	RecordID    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DroneID     string        `bson:"drone_id",json:"drone_id"`
	FaultCode   int           `bson:"fault_code",json:"fault_code"`
	Description string        `bson:"description",json:"description"`
	ReceivedOn  string        `bson:"received_on",json:"received_on"`
}

type mongoPositionRecord struct {
	RecordID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	DroneID         string        `bson:"drone_id",json:"drone_id"`
	Latitude        float32       `bson:"latitutde",json:"latitude"`
	Longitude       float32       `bson:"longitude",json:"longitude"`
	Altitude        float32       `bson:"altitude",json:"altitude"`
	CurrentSpeed    float32       `bson:"current_speed",json:"current_speed"`
	HeadingCardinal int           `bson:"heading_cardinal",json:"heading_cardinal"`
	ReceivedOn      string        `bson:"received_on",json:"received_on"`
}
