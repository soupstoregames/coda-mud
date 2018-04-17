package main

import (
	"fmt"
	"net"

	"github.com/soupstore/coda-world/config"
	"github.com/soupstore/coda-world/log"
	"github.com/soupstore/coda-world/services"
	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
	"google.golang.org/grpc"
)

func main() {
	log.Logger().Info("Starting world server")

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	sim := simulation.NewSimulation()

	// temporary
	voidID := sim.MakeRoom("Void", "Blackness. Silence. There is nothing here.")
	sim.SetSpawnRoom(voidID)
	constructID := sim.MakeRoom("The Construct", "This is the Construct. It's our loading program. We can load anything... From clothing to equipment, weapons, training simulations; anything we need.")
	sim.LinkRoom(voidID, model.East, constructID, true)
	rinseID := sim.MakeCharacter("rinse")
	sim.MakeCharacter("claw")
	sim.MakeCharacter("gesau")
	voidRoom, _ := sim.GetRoom(voidID)
	backpack := model.NewBackpack(0, "CODA Recon Pack", []string{"backpack", "pack", "recon pack", "coda recon pack", "coda recon"})
	sim.SpawnItem(backpack, voidRoom.Container.ID)
	sim.WakeUpCharacter(rinseID)
	sim.EquipItem(rinseID, backpack.ID)
	sim.SleepCharacter(rinseID)
	backpack = model.NewBackpack(0, "CODA Recon Pack", []string{"backpack", "pack", "recon pack", "coda recon pack", "coda recon"})
	sim.SpawnItem(backpack, voidRoom.Container.ID)

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
