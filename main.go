package main

import (
	"fmt"
	"github.com/soupstore/coda/components"
	"github.com/soupstore/coda/entities"
	"github.com/soupstore/coda/systems"
	"github.com/soupstore/ecs"
	"time"
)

func main() {
	world := ecs.World{}

	world.AddSystem(&systems.Weather{})
	world.AddSystem(&systems.Input{})

	CreateMap(world, 5000, 5000)

	for i := 0; i < 425088; i++ {
		CreatePlayerCharacter(world)
	}

	for {
		now := time.Now()
		world.Update(1)
		fmt.Println("Took", time.Since(now))
	}
}

func CreateMap(world ecs.World, width, height int) {
	tiles := make([][]components.Tile, width)

	for x := 0; x < width; x++ {
		tiles[x] = make([]components.Tile, height)
		for y := 0; y < height; y++ {
			tiles[x][y] = components.Tile{
				Type:      components.TileTypeOcean,
				Elevation: 0,
				Aquifer:   false,
			}
		}
	}

	gameMap := entities.GameMap{
		BasicEntity: ecs.NewBasic(),
		Geography: components.Geography{
			Width:  width,
			Height: height,
			Tiles:  tiles,
		},
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *systems.Weather:
			sys.Add(&gameMap.BasicEntity, &gameMap.Geography)
		}
	}
}

func CreatePlayerCharacter(world ecs.World) {
	player := entities.PlayerCharacter{
		BasicEntity:  ecs.NewBasic(),
		CommandQueue: components.CommandQueue{},
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *systems.Input:
			sys.Add(&player.BasicEntity, &player.CommandQueue)
		}
	}
}
