package main

import (
	"fmt"
	"net"

	"github.com/soupstore/coda-world/config"
	"github.com/soupstore/coda-world/data"
	"github.com/soupstore/coda-world/log"
	"github.com/soupstore/coda-world/services"
	"github.com/soupstore/coda-world/simulation"
	"google.golang.org/grpc"
)

func main() {
	log.Logger().Info("Starting world server")

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	sim := simulation.NewSimulation()

	// load static data
	if err := data.WatchDataFolder(conf.DataPath, sim); err != nil {
		panic(err)
	}

	// temporary
	sim.SetSpawnRoom("admin", 1)

	_ = sim.MakeCharacter("Rinse")
	sim.MakeCharacter("Claw")
	sim.MakeCharacter("Gesau")

	listenAddr := fmt.Sprintf("%s:%s", conf.Address, conf.Port)
	characterService := services.NewCharacterService(sim)

	// LISTEN TO GRPC
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Logger().Fatal(fmt.Sprintf("failed to listen: %v", err))
	}
	s := grpc.NewServer()
	services.RegisterCharacterServer(s, characterService)
	if err := s.Serve(lis); err != nil {
		log.Logger().Fatal(fmt.Sprintf("failed to serve: %v", err))
	}

}
