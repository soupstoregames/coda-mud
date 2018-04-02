package model

const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

type Direction byte

//go:generate stringer -type=Direction

func (d Direction) Opposite() Direction {
	return (d + 4) % 8
}
