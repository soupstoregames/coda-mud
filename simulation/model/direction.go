package model

import "errors"

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

func StringToDirection(d string) (Direction, error) {
	switch d {
	case "north":
		return North, nil
	case "northeast":
		return NorthEast, nil
	case "east":
		return East, nil
	case "southeast":
		return SouthEast, nil
	case "south":
		return South, nil
	case "southwest":
		return SouthWest, nil
	case "west":
		return West, nil
	case "northwest":
		return NorthWest, nil

	default:
		return North, errors.New("invalid direction")
	}
}
