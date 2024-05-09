package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Sensor struct {
	Name     string      `json:"name"`
	Location Coordinates `json:"location"`
	Tags     SensorTags  `json:"tags"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`  // between -90 and 90
	Longitude float64 `json:"longitude"` // between -180 and 180
}

type SensorTags struct {
	Unit      string `json:"unit"`
	Ingress   string `json:"ingress"`
	Distiller string `json:"distiller"`
	Name      string `json:"name"`
}

func (s *Server) ListSensors(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, sensors)
}

func (s *Server) InsertMetadata(c *gin.Context) {
	var newSensor Sensor
	if err := c.BindJSON(&newSensor); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "sensor data incorrectly formated"})
		return
	}

	sensors[newSensor.Name] = newSensor
	c.IndentedJSON(http.StatusCreated, newSensor)
}

func (s *Server) UpdateMetadata(c *gin.Context) {
	// TODO: implement
}

func (s *Server) GetMetadata(c *gin.Context) {
	name := c.Param("name")
	if sensor, ok := sensors[name]; ok {
		c.IndentedJSON(http.StatusOK, sensor)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Sensor not found in store"})
	}
}

func (s *Server) NearestLocation(c *gin.Context) {
	var userCoordinates Coordinates
	if err := c.BindJSON(&userCoordinates); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Coordinate data incorrectly formated"})
		return
	}
	// TODO: calculate nearest location
}

func (s *Server) Health(c *gin.Context) {
	if s.health {
		// It would be a good idea here to include other
		// information like time since server started and
		// server version if there was
		c.IndentedJSON(http.StatusOK, s.health)
	}
}
