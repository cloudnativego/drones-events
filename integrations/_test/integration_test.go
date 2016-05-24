package integrations_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/cf-tools"
	"github.com/cloudnativego/cfmgo"
	dronescommon "github.com/cloudnativego/drones-common"
	"github.com/cloudnativego/drones-events/mongo"
	. "github.com/cloudnativego/drones-events/service"
	"github.com/streadway/amqp"
)

var (
	appEnv, _ = cfenv.Current()
	server    = NewServer()
)

const (
	alertsQueueName    = "alerts"
	telemetryQueueName = "telemetry"
	positionsQueueName = "positions"
)

// TestIntegration submits a bunch of events on their appropriate queues, and then asserts that those events have made it to the repository in Mongo.
func TestIntegration(t *testing.T) {
	fmt.Println("Starting Integration Test.")
	ch := setupQueues(t)

	alert := dronescommon.AlertSignalledEvent{DroneID: "abc12134", Description: "bob", FaultCode: 12, ReceivedOn: time.Now().Unix()}
	dispatchMessage(ch, alert, alertsQueueName)

	position := dronescommon.PositionChangedEvent{DroneID: "drone2", Latitude: 30.12, Longitude: 25.00, Altitude: 52.00, ReceivedOn: time.Now().Unix()}
	dispatchMessage(ch, position, positionsQueueName)

	telemetry := dronescommon.TelemetryUpdatedEvent{DroneID: "drone3", CoreTemp: 25.0, RemainingBattery: 12, ReceivedOn: time.Now().Unix()}
	dispatchMessage(ch, telemetry, telemetryQueueName)

	time.Sleep(3000 * time.Millisecond)

	rollupRepo := initRepository(t)
	alertEvent, err := rollupRepo.GetAlertEvent("abc12134")
	if err != nil {
		t.Errorf("Failed checking repo for dispatched alert: %s\n", err.Error())
		return
	}
	if alertEvent.FaultCode != 12 {
		t.Errorf("Alert event in repo doesn't match dispatched event: orig: %+v, repo: %+v\n", alert, alertEvent)
		return
	}

	positionEvent, err := rollupRepo.GetPositionEvent("drone2")
	if err != nil {
		t.Errorf("Failed checking repo for dispatched position event: %s\n", err.Error())
		return
	}
	if positionEvent.Altitude != 52.00 {
		t.Errorf("Position event in repo doesn't match dispatched event: orig: %+v, repo: %+v\n", position, positionEvent)
		return
	}

	telemetryEvent, err := rollupRepo.GetTelemetryEvent("drone3")
	if err != nil {
		t.Errorf("Failed checking repo for dispatched telemetry event: %s\n", err.Error())
		return
	}
	if telemetryEvent.CoreTemp != 25.0 {
		t.Errorf("Telemetry event in repo doesn't match dispatched event: orig: %+v, repo: %+v\n", telemetry, telemetryEvent)
		return
	}
}

/*
 * === Utility Functions ===
 */

func setupQueues(t *testing.T) (channel *amqp.Channel) {
	amqpURI, err := cftools.GetVCAPServiceProperty("rabbit", "url", appEnv)
	if err != nil {
		t.Errorf("No AMQP connection detected... Integration FAIL.\n")
		return
	}
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		t.Errorf("Failed to connect to rabbit, %v\n", err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		t.Errorf("Failed to open AMQP channel %v\n", err)
		return
	}

	ch.QueueDeclare(
		alertsQueueName, // name
		false,           // durable
		false,           // delete when usused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	ch.QueueDeclare(
		positionsQueueName, // name
		false,              // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	ch.QueueDeclare(
		telemetryQueueName, // name
		false,              // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	channel = ch
	return
}

func dispatchMessage(c *amqp.Channel, message interface{}, queueName string) (err error) {
	body, err := json.Marshal(message)
	if err == nil {
		err = c.Publish(
			"",        // exchange
			queueName, // routing key
			true,      // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			fmt.Printf("Failed to dispatch message: %s\n", err)
		}
	} else {
		fmt.Printf("Failed to marshal message %v (%s)\n", message, err)
	}
	return
}

func initRepository(t *testing.T) (repo eventRepository) {
	dbServiceURI, err := cftools.GetVCAPServiceProperty("mongo-eventrollup", "url", appEnv)
	if err != nil || len(dbServiceURI) == 0 {
		t.Errorf("\nError retreieving database configuration: %v\n", err)
	} else {
		telemetryCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, "telemetry")
		positionsCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, "positions")
		alertsCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, "alerts")
		repo = mongo.NewEventRollupRepository(positionsCollection, alertsCollection, telemetryCollection)
	}
	return
}

type eventRepository interface {
	GetTelemetryEvent(droneID string) (event dronescommon.TelemetryUpdatedEvent, err error)
	GetPositionEvent(droneID string) (event dronescommon.PositionChangedEvent, err error)
	GetAlertEvent(droneID string) (event dronescommon.AlertSignalledEvent, err error)
}
