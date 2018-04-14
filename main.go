package main

import (
	"fmt"
	"log"
	"net"

	"github.com/soupstore/coda-world/config"
	"github.com/soupstore/coda-world/services"
	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	//logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	sim := simulation.NewSimulation(logger)

	// temporary
	voidID := sim.MakeRoom("Void", "Blackness. Silence. There is nothing here.")
	sim.SetSpawnRoom(voidID)
	constructID := sim.MakeRoom("The Construct", "This is the Construct. It's our loading program. We can load anything... From clothing to equipment, weapons, training simulations; anything we need.")
	sim.LinkRoom(voidID, model.East, constructID, true)
	sim.MakeCharacter("rinse")
	sim.MakeCharacter("claw")
	sim.MakeCharacter("gesau")
	voidRoom, _ := sim.GetRoom(voidID)
	sim.SpawnItem(model.NewBackpack(0, "CODA Recon Pack"), voidRoom.Container.ID)

	listenAddr := fmt.Sprintf("%s:%s", conf.Address, conf.Port)
	characterService := services.NewCharacterService(sim, logger)

	// LISTEN TO GRPC
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	services.RegisterCharacterServer(s, characterService)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
