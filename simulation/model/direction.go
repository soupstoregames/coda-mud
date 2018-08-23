package model

import "errors"

const (
	DirectionNorth Direction = iota
	DirectionNorthEast
	DirectionEast
	DirectionSouthEast
	DirectionSouth
	DirectionSouthWest
	DirectionWest
	DirectionNorthWest
)

// Direction is an enum of 8-point compass directions.
type Direction byte

// Opposite returns the direction that faces the other way. North would give south.
func (d Direction) Opposite() Direction {
	return (d + 4) % 8
}

func (d Direction) String() string {
	switch d {
	case DirectionNorth:
		return "north"
	case DirectionNorthEast:
		return "northeast"
	case DirectionEast:
		return "east"
	case DirectionSouthEast:
		return "southeast"
	case DirectionSouth:
		return "south"
	case DirectionSouthWest:
		return "southwest"
	case DirectionWest:
		return "west"
	case DirectionNorthWest:
		return "northwest"

	default:
		return "Invalid direction"
	}
}

// StringToDirection attempts to parse a string into a Direction.
// If unable, it returns North and an error.
func StringToDirection(d string) (Direction, error) {
	switch d {
	case "north":
		return DirectionNorth, nil
	case "northeast":
		return DirectionNorthEast, nil
	case "east":
		return DirectionEast, nil
	case "southeast":
		return DirectionSouthEast, nil
	case "south":
		return DirectionSouth, nil
	case "southwest":
		return DirectionSouthWest, nil
	case "west":
		return DirectionWest, nil
	case "northwest":
		return DirectionNorthWest, nil

	default:
		return DirectionNorth, errors.New("invalid direction")
	}
}
