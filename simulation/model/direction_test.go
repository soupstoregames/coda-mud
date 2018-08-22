package model_test

import (
	"testing"

	"github.com/soupstore/coda/simulation/model"
)

func TestOpposite(t *testing.T) {
	tests := []struct {
		s model.Direction
		n model.Direction
	}{
		{model.DirectionNorth, model.DirectionSouth},
		{model.DirectionNorthEast, model.DirectionSouthWest},
		{model.DirectionEast, model.DirectionWest},
		{model.DirectionSouthEast, model.DirectionNorthWest},
		{model.DirectionSouth, model.DirectionNorth},
		{model.DirectionSouthWest, model.DirectionNorthEast},
		{model.DirectionWest, model.DirectionEast},
		{model.DirectionNorthWest, model.DirectionSouthEast},
	}

	for _, test := range tests {
		expected := test.n
		actual := test.s.Opposite()

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		d model.Direction
		s string
	}{
		{model.DirectionNorth, "DirectionNorth"},
		{model.DirectionNorthEast, "DirectionNorthEast"},
		{model.DirectionEast, "DirectionEast"},
		{model.DirectionSouthEast, "DirectionSouthEast"},
		{model.DirectionSouth, "DirectionSouth"},
		{model.DirectionSouthWest, "DirectionSouthWest"},
		{model.DirectionWest, "DirectionWest"},
		{model.DirectionNorthWest, "DirectionNorthWest"},
		{9, "Direction(9)"},
	}

	for _, test := range tests {
		expected := test.s
		actual := test.d.String()

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}
	}
}
