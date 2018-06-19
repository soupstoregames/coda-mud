package main

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/soupstore/coda/common/config"
	"github.com/soupstore/coda/common/logging"
	"github.com/soupstore/coda/database"
	"github.com/soupstore/coda/database/migrations"
	"github.com/soupstore/coda/simulation"
	static "github.com/soupstore/coda/static-data"
	"github.com/soupstore/coda/telnet"
)

func main() {
	logging.Logger().Info("Starting world server")

	// load config values from env vars
	conf, err := config.Load()
	if err != nil {
		logging.Logger().Fatal(err.Error())
	}

	// connect to DB
	db := pg.Connect(&pg.Options{
		User:     "Nick",
		Password: "",
		Database: "coda",
	})

	// run database migrations
	migrationAsset := database.MakeBinDataMigration(migrations.AssetNames(), migrations.Asset)
	err = database.PerformMigration(migrationAsset)
	if err != nil {
		log.Fatal(err.Error())
	}

	// create a new simulation
	sim := simulation.NewSimulation(db)

	// load static data and apply changes to simulation
	dw := static.NewDataWatcher(conf.DataPath, sim)
	logging.SubscribeToErrorChan(dw.Errors)

	// load state data and apply to the simulation
	characters, err := database.GetCharacters(db)
	if err != nil {
		logging.Logger().Fatal(err.Error())
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
