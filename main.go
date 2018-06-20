package main

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/soupstore/coda/common/config"
	"github.com/soupstore/coda/common/logging"
	"github.com/soupstore/coda/database"
	"github.com/soupstore/coda/database/migrations"
	"github.com/soupstore/coda/services"
	"github.com/soupstore/coda/simulation"
	static "github.com/soupstore/coda/static-data"
	"github.com/soupstore/coda/telnet"
)

func main() {
	var (
		conf         *config.Config
		db           *pg.DB
		sim          *simulation.Simulation
		usersManager *services.UsersManager

		err error
	)

	logging.Logger().Info("Starting world server")

	if conf, err = config.Load(); err != nil {
		logging.Logger().Fatal(err.Error())
	}

	if db, err = connectToDatabaseAndMigrate(conf); err != nil {
		logging.Logger().Fatal(err.Error())
	}

	if sim, err = createAndInitializeSimulation(conf, db); err != nil {
		logging.Logger().Fatal(err.Error())
	}

	usersManager = services.NewUsersManagers(db)

	// temporary
	sim.SetSpawnRoom("arrival-city", 1)

	err = launchTelnetServer(conf, sim, usersManager)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func connectToDatabaseAndMigrate(conf *config.Config) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     conf.DatabaseUser,
		Password: conf.DatabasePassword,
		Database: conf.DatabaseName,
	})

	migrationAsset := database.MakeBinDataMigration(migrations.AssetNames(), migrations.Asset)
	err := database.PerformMigration(migrationAsset)

	return db, err
}

func createAndInitializeSimulation(conf *config.Config, db *pg.DB) (*simulation.Simulation, error) {
	// create a new simulation
	sim := simulation.NewSimulation(db)

	// load static data and apply changes to simulation
	dw := static.NewDataWatcher(conf.DataPath, sim)
	logging.SubscribeToErrorChan(dw.Errors)

	// load state data and apply to the simulation
	if err := loadSavedState(db, sim); err != nil {
		return nil, err
	}

	return sim, nil
}

func loadSavedState(db *pg.DB, sim *simulation.Simulation) error {
	characters, err := database.GetCharacters(db)
	if err != nil {
		return err
	}

	sim.LoadCharacters(characters)
	return nil
}

func launchTelnetServer(conf *config.Config, sim *simulation.Simulation, usersManager *services.UsersManager) error {
	listenAddr := fmt.Sprintf("%s:%s", conf.Address, conf.Port)
	telnetServer := telnet.NewServer(conf, listenAddr, sim, usersManager)
	return telnetServer.ListenAndServe()
}
