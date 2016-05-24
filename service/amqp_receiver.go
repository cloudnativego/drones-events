package service

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/cf-tools"
	common "github.com/cloudnativego/drones-common"
	"github.com/streadway/amqp"
)

func dequeueEvents(alertChannel chan common.AlertSignalledEvent, telemetryChannel chan common.TelemetryUpdatedEvent, positionChannel chan common.PositionChangedEvent) {
	fmt.Printf("Starting AMQP queue de-serializer...")
	appEnv, _ := cfenv.Current()
	amqpURI, err := cftools.GetVCAPServiceProperty("rabbit", "url", appEnv)
	if err != nil {
		fmt.Println("No Rabbit/AMQP connection details supplied. ABORTING. No events will be dequeued!!!")
		return
	}
	fmt.Printf("dialing %s\n", amqpURI)
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		fmt.Printf("Failed to connect to rabbit, %v\n", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open AMQP channel %v\n", err)
	}

	alertsQ, _ := ch.QueueDeclare(
		alertsQueueName, // name
		false,           // durable
		false,           // delete when usused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	positionsQ, _ := ch.QueueDeclare(
		positionsQueueName, // name
		false,              // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	telemetryQ, _ := ch.QueueDeclare(
		telemetryQueueName, // name
		false,              // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	alertsIn, _ := ch.Consume(
		alertsQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	positionsIn, _ := ch.Consume(
		positionsQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	telemetryIn, _ := ch.Consume(
		telemetryQ.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for {
			select {
			case alertRaw := <-alertsIn:
				dispatchAlert(alertRaw, alertChannel)
			case telemetryRaw := <-telemetryIn:
				dispatchTelemetry(telemetryRaw, telemetryChannel)
			case positionRaw := <-positionsIn:
				dispatchPosition(positionRaw, positionChannel)
			}
		}
	}()
}

func dispatchAlert(alertRaw amqp.Delivery, out chan common.AlertSignalledEvent) {
	var event common.AlertSignalledEvent
	err := json.Unmarshal(alertRaw.Body, &event)
	if err == nil {
		out <- event
	} else {
		fmt.Printf("Failed to de-serialize raw alert from queue, %v\n", err)
	}
	return
}

func dispatchTelemetry(telemetryRaw amqp.Delivery, out chan common.TelemetryUpdatedEvent) {
	var event common.TelemetryUpdatedEvent
	err := json.Unmarshal(telemetryRaw.Body, &event)
	if err == nil {
		out <- event
	} else {
		fmt.Printf("Failed to de-serialize raw telemetry from queue, %v\n", err)
	}
	return
}

func dispatchPosition(positionRaw amqp.Delivery, out chan common.PositionChangedEvent) {
	var event common.PositionChangedEvent
	err := json.Unmarshal(positionRaw.Body, &event)
	if err == nil {
		out <- event
	} else {
		fmt.Printf("Failed to de-serialize raw position from queue, %v\n", err)
	}
	return
}
