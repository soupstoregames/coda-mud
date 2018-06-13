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
		{model.North, model.South},
		{model.NorthEast, model.SouthWest},
		{model.East, model.West},
		{model.SouthEast, model.NorthWest},
		{model.South, model.North},
		{model.SouthWest, model.NorthEast},
		{model.West, model.East},
		{model.NorthWest, model.SouthEast},
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
		{model.North, "North"},
		{model.NorthEast, "NorthEast"},
		{model.East, "East"},
		{model.SouthEast, "SouthEast"},
		{model.South, "South"},
		{model.SouthWest, "SouthWest"},
		{model.West, "West"},
		{model.NorthWest, "NorthWest"},
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
