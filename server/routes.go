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
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
	}
	c.IndentedJSON(http.StatusOK, sensors)
}

func (s *Server) addSensor(c *gin.Context) {
	var newSensor *Sensor
	if err := c.BindJSON(&newSensor); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := s.insertSensor(
		newSensor.Name,
		newSensor.Tags.Unit,
		newSensor.Location.Latitude,
		newSensor.Location.Longitude,
	); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
	}
	c.IndentedJSON(http.StatusCreated, newSensor)
}

func (s *Server) getSensor(c *gin.Context) {
	var err error
	var sensor *Sensor
	if sensor, err = s.querySensor(c.Param("name")); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.IndentedJSON(http.StatusOK, sensor)
}

func (s *Server) nearestSensor(c *gin.Context) {
	// Parse the query parameters to float64
	var err error
	var latitude, longitude float64
	if latitude, err = strconv.ParseFloat(c.Param("lat"), 64); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
	}
	if longitude, err = strconv.ParseFloat(c.Param("lon"), 64); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
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

type Status struct {
	Ok     bool   `json:"ok"`
	Uptime string `json:"uptime"`
}

// Health check for server. Usually It would we
// should include the server version if there was one.
func (s *Server) statusCheck(c *gin.Context) {
	status := Status{
		Ok:     s.healthy,
		Uptime: time.Since(s.started).String(),
	}
	c.IndentedJSON(http.StatusOK, status)
}
