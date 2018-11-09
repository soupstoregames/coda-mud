package model

import "errors"

const (
	DirectionNorth Direction = iota
	DirectionNorthEast
	DirectionEast
	DirectionSouthEast
	DirectionUp
	DirectionSouth
	DirectionSouthWest
	DirectionWest
	DirectionNorthWest
	DirectionDown
)

// Direction is an enum of 8-point compass directions.
type Direction byte

// Opposite returns the direction that faces the other way. North would give south.
func (d Direction) Opposite() Direction {
	return (d + 5) % 10
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
	case DirectionUp:
		return "up"
	case DirectionSouth:
		return "south"
	case DirectionSouthWest:
		return "southwest"
	case DirectionWest:
		return "west"
	case DirectionNorthWest:
		return "northwest"
	case DirectionDown:
		return "down"

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
	case "up":
		return DirectionUp, nil
	case "south":
		return DirectionSouth, nil
	case "southwest":
		return DirectionSouthWest, nil
	case "west":
		return DirectionWest, nil
	case "northwest":
		return DirectionNorthWest, nil
	case "down":
		return DirectionDown, nil

	default:
		return DirectionNorth, errors.New("invalid direction")
	}
}
