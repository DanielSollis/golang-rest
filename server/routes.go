package server

import (
	"math"
	"net/http"
	"time"

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

type Status struct {
	Ok     bool   `json:"ok"`
	Uptime string `json:"uptime"`
}

// TODO
func (s *Server) ListSensors(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, sensors)
}

// TODO
func (s *Server) InsertSensor(c *gin.Context) {
	var newSensor Sensor
	if err := c.BindJSON(&newSensor); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "sensor data incorrectly formated"})
		return
	}

	sensors[newSensor.Name] = newSensor
	c.IndentedJSON(http.StatusCreated, newSensor)
}

// TODO
func (s *Server) UpdateSensor(c *gin.Context) {
	// TODO: implement
}

// TODO
func (s *Server) GetSensor(c *gin.Context) {
	name := c.Param("name")
	if sensor, ok := sensors[name]; ok {
		c.IndentedJSON(http.StatusOK, sensor)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Sensor not found in store"})
	}
}

// TODO
func (s *Server) NearestLocation(c *gin.Context) {
	var userCoordinates Coordinates
	if err := c.BindJSON(&userCoordinates); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "coordinate data incorrectly formated"})
		return
	}

	latitude := userCoordinates.Latitude
	if latitude < -90 || latitude > 90 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "latitude must be between -90 and 90"})
		return
	}

	longitude := userCoordinates.Longitude
	if longitude < -180 || longitude > 180 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "longitude must be between -180 and 180"})
		return
	}

	min := math.Inf(1)
	var minSensor Sensor
	for _, sensor := range sensors {
		distance := Haversine(userCoordinates, sensor.Location)
		min = math.Min(min, distance)
		minSensor = sensor
	}
	c.IndentedJSON(http.StatusOK, minSensor)
}

// Calculates the great circle distance between two points.
// The haversine formula assumes points on a perfect sphere
// (the earth isn't a perfect sphere) so the haversine error
// can be up to 0.5%
func Haversine(user, sensor Coordinates) float64 {
	userLat, sensorLat := user.Latitude, sensor.Latitude
	userLong, sensorLong := user.Longitude, sensor.Longitude

	// Distance between latitudes in radians
	latDistanceRad := (sensorLat - userLat) * math.Pi / 180
	longDistanceRad := (sensorLong - userLong) * math.Pi / 180

	// Latitudes in radians
	userLatRad := userLat * math.Pi / 180
	sensorLatRad := sensorLat * math.Pi / 180

	// Calculate the square of half the chord length between two points 'a'
	latPower := math.Pow(math.Sin(latDistanceRad/2), 2)
	longPower := math.Pow(math.Sin(longDistanceRad/2), 2)
	latCosine := math.Cos(userLatRad) * math.Cos(sensorLatRad)
	a := (latPower + longPower) * latCosine

	// Calculate the angular between the two points 'c'
	c := 2 * math.Asin(math.Sqrt(a))

	earthsRadius := 6371
	return float64(earthsRadius) * c
}

// Health check for server. Usually It would be
// best practice to include other info here like time
// since the server started and server version if
// there was one.
func (s *Server) Health(c *gin.Context) {
	status := Status{
		Ok:     s.health,
		Uptime: time.Since(s.started).String(),
	}
	c.IndentedJSON(http.StatusOK, status)
}
