package components

type TileType string

const (
	TileTypePlains   TileType = "plains"
	TileTypeForest   TileType = "forest"
	TileTypeRiver    TileType = "river"
	TileTypeOcean    TileType = "ocean"
	TileTypeRocky    TileType = "rocky"
	TileTypeMountain TileType = "mountain"
	TileTypePaved    TileType = "paved"
	TileTypeRoad     TileType = "road"
)

type Tile struct {
	Elevation byte
	Type      TileType
	Aquifer   bool
}

type Geography struct {
	Width  int
	Height int
	Tiles  [][]Tile
}
