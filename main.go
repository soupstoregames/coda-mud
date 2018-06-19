package main

import (
	"fmt"

	"github.com/soupstore/coda/common/config"
	"github.com/soupstore/coda/common/log"
	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/coda/simulation/data/state"
	"github.com/soupstore/coda/simulation/data/static"
	"github.com/soupstore/coda/telnet"
)

func main() {
	log.Logger().Info("Starting world server")

	// load config values from env vars
	conf, err := config.Load()
	if err != nil {
		log.Logger().Fatal(err.Error())
	}

	db, err := state.OpenConnection("Nick", "", "coda", "")
	if err != nil {
		log.Logger().Fatal(err.Error())
	}

	// create a new simulation
	sim := simulation.NewSimulation(db)

	// load static data and apply changes to simulation
	dw := static.NewDataWatcher(conf.DataPath, sim)
	log.SubscribeToErrorChan(dw.Errors)

	// load state data and apply to the simulation
	characters, err := state.GetCharacters(db)
	if err != nil {
		log.Logger().Fatal(err.Error())
	}
	sim.LoadCharacters(characters)

	// temporary
	sim.SetSpawnRoom("arrival-city", 1)

	// _ = sim.MakeCharacter("Rinse")
	// sim.MakeCharacter("Claw")
	// sim.MakeCharacter("Gesau")

	spawnRoom, _ := sim.GetRoom("arrival-city", 1)
	sim.SpawnItem(1, spawnRoom.Container.ID)
	sim.SpawnItem(2, spawnRoom.Container.ID)

	listenAddr := fmt.Sprintf("%s:%s", conf.Address, conf.Port)
	telnetServer := telnet.NewServer(conf, listenAddr, sim)
	err = telnetServer.ListenAndServe()
}
