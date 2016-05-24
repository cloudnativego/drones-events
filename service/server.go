package service

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/cf-tools"
	"github.com/cloudnativego/cfmgo"
	dronescommon "github.com/cloudnativego/drones-common"
	"github.com/cloudnativego/drones-events/mongo"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)

	alertChannel := make(chan dronescommon.AlertSignalledEvent)
	telemetryChannel := make(chan dronescommon.TelemetryUpdatedEvent)
	positionChannel := make(chan dronescommon.PositionChangedEvent)

	repo := InitRepository()
	dequeueEvents(alertChannel, telemetryChannel, positionChannel)
	consumeEvents(alertChannel, telemetryChannel, positionChannel, repo)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/", homeHandler(formatter)).Methods("GET")
}

// InitRepository attempts to create a mongo repository
func InitRepository() (repo eventRepository) {
	appEnv, _ := cfenv.Current()
	dbServiceURI, err := cftools.GetVCAPServiceProperty("mongo-eventrollup", "url", appEnv)
	if err != nil || len(dbServiceURI) == 0 {
		if err != nil {
			fmt.Printf("\nError retreieving database configuration: %v\n", err)
		}
		fmt.Println("MongoDB was not detected, using fake repository (THIS IS BAD)...")
		repo = newFakeRepository()
	} else {
		telemetryCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, "telemetry")
		positionsCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, "positions")
		alertsCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, "alerts")
		repo = mongo.NewEventRollupRepository(positionsCollection, alertsCollection, telemetryCollection)
	}
	return
}
