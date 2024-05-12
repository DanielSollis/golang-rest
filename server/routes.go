package server

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) listSensors(c *gin.Context) {
	sensors, err := s.queryAllSensors()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, sensors)
}

func (s *Server) addSensor(c *gin.Context) {
	var newSensor *Sensor
	if err := c.BindJSON(&newSensor); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.insertSensor(
		newSensor.Name,
		newSensor.Tags.Unit,
		newSensor.Location.Latitude,
		newSensor.Location.Longitude,
	); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusCreated, newSensor)
}

func (s *Server) updateSensor(c *gin.Context) {
	var toUpdate *Sensor
	if err := c.BindJSON(&toUpdate); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.updateSensorStore(
		toUpdate.Name,
		toUpdate.Tags.Unit,
		toUpdate.Location.Latitude,
		toUpdate.Location.Longitude,
	); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, toUpdate)
}

func (s *Server) getSensor(c *gin.Context) {
	var err error
	var sensor *Sensor
	if sensor, err = s.querySensor(c.Param("name")); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, sensor)
}

func (s *Server) getNearestSensor(c *gin.Context) {
	// Parse the query parameters to float64
	var err error
	var latitude, longitude float64
	if latitude, err = strconv.ParseFloat(c.Param("lat"), 64); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if longitude, err = strconv.ParseFloat(c.Param("lon"), 64); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Input validation
	if latitude < -90 || latitude > 90 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "latitude must be between -90 and 90"})
		return
	}
	if longitude < -180 || longitude > 180 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "longitude must be between -180 and 180"})
		return
	}

	// Query all sensors
	sensors, err := s.queryAllSensors()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	// Find the nearest sensor
	min := math.Inf(1)
	var minSensor *Sensor
	userCoordinates := &Coordinates{Latitude: latitude, Longitude: longitude}
	for _, sensor := range sensors {
		distance := haversine(userCoordinates, &sensor.Location)
		if distance < min {
			min = distance
			minSensor = sensor
		}
	}
	c.IndentedJSON(http.StatusOK, minSensor)
}

// Calculates the great circle distance between two points.
// The haversine formula assumes points on a perfect sphere
// (the earth isn't a perfect sphere) so the haversine error
// can be up to 0.5%
func haversine(user, sensor *Coordinates) float64 {
	userLat, sensorLat := user.Latitude, sensor.Latitude
	userLon, sensorLon := user.Longitude, sensor.Longitude

	// Distance between latitudes in radians
	latDistanceRad := (sensorLat - userLat) * math.Pi / 180
	lonDistanceRad := (sensorLon - userLon) * math.Pi / 180

	// Latitudes in radians
	userLatRad := userLat * math.Pi / 180
	sensorLatRad := sensorLat * math.Pi / 180

	// Calculate the square of half the chord length between two points 'a'
	latPower := math.Pow(math.Sin(latDistanceRad/2), 2)
	lonPower := math.Pow(math.Sin(lonDistanceRad/2), 2)
	latCosine := math.Cos(userLatRad) * math.Cos(sensorLatRad)
	a := latPower + lonPower*latCosine

	// Calculate the angular between the two points 'c'
	c := 2 * math.Asin(math.Sqrt(a))

	// In km (...approximately)
	earthsRadius := 6371

	return float64(earthsRadius) * c
}

type Status struct {
	Ok     bool   `json:"ok"`
	Uptime string `json:"uptime"`
}

// Health check for server. Usually It would we
// should include the server version if there was one
func (s *Server) statusCheck(c *gin.Context) {
	status := Status{
		Ok:     s.healthy,
		Uptime: time.Since(s.started).String(),
	}
	c.IndentedJSON(http.StatusOK, status)
}
