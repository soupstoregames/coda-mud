package main

import (
	"fmt"
	"log"
	"time"

	"github.com/soupstore/coda/config"
	"github.com/soupstore/coda/data/state"
	"github.com/soupstore/coda/data/static"
	"github.com/soupstore/coda/logging"
	"github.com/soupstore/coda/servers/telnet"
	"github.com/soupstore/coda/services"
	"github.com/soupstore/coda/simulation"
)

func main() {
	var (
		conf         *config.Config
		persister    state.Persister
		sim          *simulation.Simulation
		usersManager *services.UsersManager
		err          error
	)

	logging.Logger().Info("Starting world server")

	if conf, err = config.Load(); err != nil {
		logging.Logger().Fatal(err.Error())
	}

	if persister, err = state.NewFileSystemPersister(conf); err != nil {
		logging.Logger().Fatal(err.Error())
	}

	if sim, err = createAndInitializeSimulation(conf, persister); err != nil {
		logging.Logger().Fatal(err.Error())
	}

	go func() {
		for {
			time.Sleep(3 * time.Second)
			if err := sim.Save(); err != nil {
				panic(err)
			}
		}
	}()

	usersManager = services.NewUsersManager()

	// temporary
	sim.SetSpawnRoom("arrival-city", 1)
	room, err := sim.GetRoom("arrival-city", 1)
	if err != nil {
		panic(err)
	}
	usersManager.Register("rinse", "bums")
	id := sim.MakeCharacter("Rinse")
	if err := sim.SpawnItem(1, room.Container.ID()); err != nil {
		panic(err)
	}

	usersManager.AssociateCharacter("rinse", id)

	err = launchTelnetServer(conf, sim, usersManager)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func createAndInitializeSimulation(conf *config.Config, persister state.Persister) (sim *simulation.Simulation, err error) {
	sim = simulation.NewSimulation(persister)
	dw := static.NewDataWatcher(conf.DataPath, sim)
	logging.SubscribeToErrorChan(dw.Errors)
	return
}

func launchTelnetServer(conf *config.Config, sim *simulation.Simulation, usersManager *services.UsersManager) error {
	listenAddr := fmt.Sprintf("%s:%s", conf.Address, conf.Port)
	telnetServer := telnet.NewServer(conf, listenAddr, sim, usersManager)
	return telnetServer.ListenAndServe()
}
