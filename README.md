# Drone Army Event Processor
This is a microservice responsible for event processing in the Drone Army Event Sourcing / CQRS example.

It will bring up a service that allows it to remain up at all times while it reads incoming events from the
queue and then updates the aggregate state of the system accordingly.
