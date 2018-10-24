package main

import (
	"fmt"
	"log"
	"time"

	"github.com/soupstore/coda/config"
	"github.com/soupstore/coda/servers/telnet"
	"github.com/soupstore/coda/services"
	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/coda/simulation/data/state"
	"github.com/soupstore/coda/simulation/data/static"
	"github.com/soupstore/go-core/logging"
)

func main() {
	var (
		conf         *config.Config
		staticData   *static.DataWatcher
		stateData    *state.FileSystem
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

	// create a persister to save the simulation state
	if stateData, err = state.NewFileSystem(conf); err != nil {
		logging.Fatal(err.Error())
	}

	// create the users service for managing login details
	usersManager = services.NewUsersManager()

	// load the static data
	if err := staticData.InitialLoad(); err != nil {
		logging.Fatal(err.Error())
	}

	// start watching for changes to the static data folder
	staticData.Watch()

	// load the saved state
	loadState(stateData, usersManager, sim)

	// temporary
	sim.SetSpawnRoom("arrival-city", 1)
	// room, err := sim.GetRoom("arrival-city", 1)
	// if err != nil {
	// 	logging.Fatal(err.Error())
	// }
	// usersManager.Register("rinse", "bums")
	// id := sim.MakeCharacter("Rinse")
	// usersManager.AssociateCharacter("rinse", id)

	// set up save timing for simulation state
	startSaveSimulationTicker(usersManager, sim, stateData)

	// start the simulation

	// start the telnet server
	telnetServer := telnet.NewServer(conf, sim, usersManager)
	if err = telnetServer.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

}

func loadState(loader state.Loader, usersManager *services.UsersManager, sim simulation.StateController) {
	var (
		users      []state.User
		characters []state.Character
		worlds     []state.World
		err        error
	)

	if users, characters, worlds, err = loader.Load(); err != nil {
		logging.Fatal(err.Error())
	}
	if usersManager.Load(users); err != nil {
		logging.Fatal(err.Error())
	}
	if sim.Load(characters, worlds); err != nil {
		logging.Fatal(err.Error())
	}
}

func startSaveSimulationTicker(u *services.UsersManager, s *simulation.Simulation, p state.Persister) {
	t := time.NewTicker(time.Minute)
	go func() {
		for range t.C {
			if err := u.Save(p); err != nil {
				logging.Warn(fmt.Sprintf("Failed to save users: %s", err.Error()))
			}
			if err := s.Save(p); err != nil {
				logging.Warn(fmt.Sprintf("Failed to save simulation state: %s", err.Error()))
			}
		}
	}()
}
