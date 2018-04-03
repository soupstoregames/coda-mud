package main

import (
	"fmt"
	"log"
	"net"

	"github.com/soupstore/coda-world/config"
	"github.com/soupstore/coda-world/services"
	"github.com/soupstore/coda-world/simulation"
	"google.golang.org/grpc"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	sim := simulation.NewSimulation()

	listenAddr := fmt.Sprintf("%s:%s", conf.Address, conf.Port)
	characterService := services.NewCharacterService(sim)

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
