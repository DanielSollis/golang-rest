package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	sync.RWMutex
	srv     *http.Server // http server for API defaults
	router  *gin.Engine  // the http handler
	health  bool         // server state for health checks
	errChan chan error   // synchronize gracefull shutdown
}

var sensors = map[string]Sensor{
	"1": {Name: "1", Location: Coordinates{}, Tags: SensorTags{}},
	"2": {Name: "2", Location: Coordinates{}, Tags: SensorTags{}},
	"3": {Name: "3", Location: Coordinates{}, Tags: SensorTags{}},
}

// TODO
func (s *Server) serve() {
	s.setStatus(true)
}

// TODO
func (s *Server) shutdown() {
	// if this is more complicated, return errors, lock, etc.

	log.Println("Gracefully shutting down server")
	s.setStatus(false)
	log.Println("Server successfully shutdown")
}

// TODO
func (s *Server) setStatus(health bool) {
	s.Lock()
	s.health = health
	s.Unlock()
}

// TODO
func (s *Server) setupRoutes() {
	s.router.GET("/allsensors", s.ListSensors)
	s.router.POST("/sensor", s.InsertMetadata)
	s.router.PUT("/sensor", s.UpdateMetadata)
	s.router.GET("/sensor:name", s.GetMetadata)
	s.router.GET("/nearest", s.NearestLocation)
	s.router.GET("/health", s.Health)
}

// TODO
func main() {

}
