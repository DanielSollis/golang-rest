package server_test

import (
	"math"
	"pingthings/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListSensors(t *testing.T) {
	// TODO
}

func TestAddSensor(t *testing.T) {
	// TODO
}

func TestGetSensor(t *testing.T) {
	// TODO
}

func TestNearestSensor(t *testing.T) {
	// TODO
}

func TestStatusCheck(t *testing.T) {

}

func TestHaversine(t *testing.T) {
	// test one
	userCoordinates := server.Coordinates{
		Latitude:  0,
		Longitude: 0,
	}
	sensorCoordinates := server.Coordinates{
		Latitude:  0,
		Longitude: 180,
	}
	distance := server.Haversine(userCoordinates, sensorCoordinates)
	expected := 20015
	assert.EqualValues(t, expected, math.Round(distance))

	// Test two
	userCoordinates = server.Coordinates{
		Latitude:  51.5007,
		Longitude: 0.1246,
	}
	sensorCoordinates = server.Coordinates{
		Latitude:  40.6892,
		Longitude: 74.0445,
	}
	distance = server.Haversine(userCoordinates, sensorCoordinates)
	expected = 5575
	assert.EqualValues(t, expected, math.Round(distance))
}
