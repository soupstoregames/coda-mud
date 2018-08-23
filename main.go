package main

import (
	"fmt"
	"log"
	"time"

	"github.com/soupstore/coda/config"
	"github.com/soupstore/coda/data/state"
	"github.com/soupstore/coda/data/static"
	"github.com/soupstore/coda/servers/telnet"
	"github.com/soupstore/coda/services"
	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/go-core/logging"
)

func main() {
	var (
		conf         *config.Config
		staticData   *static.DataWatcher
		stateData    state.Persister
		sim          *simulation.Simulation
		usersManager *services.UsersManager
		err          error
	)

	logging.Info("Starting")

	// load the configuration from environmental variables
	if conf, err = config.Load(); err != nil {
		logging.Fatal(err.Error())
	}

	// create the simulation
	sim = simulation.NewSimulation()

	// create the static data loader
	staticData = static.NewDataWatcher(conf.DataPath, sim)
	logging.SubscribeToErrorChan(staticData.Errors)

	// load the static data
	if err := staticData.InitialLoad(); err != nil {
		logging.Fatal(err.Error())
	}

	// create a persister to save the simulation state
	if stateData, err = state.NewFileSystemPersister(conf); err != nil {
		logging.Fatal(err.Error())
	}

	// load the saved state

	// create the users service for managing login details
	usersManager = services.NewUsersManager()

	// temporary
	sim.SetSpawnRoom("arrival-city", 1)
	room, err := sim.GetRoom("arrival-city", 1)
	if err != nil {
		logging.Fatal(err.Error())
	}
	usersManager.Register("rinse", "bums")
	id := sim.MakeCharacter("Rinse")
	if err := sim.SpawnItem(1, room.Container.ID()); err != nil {
		logging.Fatal(err.Error())
	}
	if err := sim.SpawnItem(2, room.Container.ID()); err != nil {
		logging.Fatal(err.Error())
	}
	if err := sim.SpawnItem(2, room.Container.ID()); err != nil {
		logging.Fatal(err.Error())
	}
	usersManager.AssociateCharacter("rinse", id)

	// start watching for changes to the static data folder
	staticData.Watch()

	// set up save timing for simulation state
	startSaveSimulationTicker(sim, stateData)

	// start the simulation

	// start the telnet server
	telnetServer := telnet.NewServer(conf, sim, usersManager)
	if err = telnetServer.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

}

func startSaveSimulationTicker(s *simulation.Simulation, p state.Persister) {
	t := time.NewTicker(time.Minute)
	go func() {
		for range t.C {
			if err := s.Save(p); err != nil {
				logging.Warn(fmt.Sprintf("Failed to save simulation state: %s", err.Error()))
			}
		}
	}()
}
