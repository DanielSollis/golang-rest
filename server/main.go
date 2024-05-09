package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	sync.RWMutex
	srv     *http.Server // http server for API defaults
	router  *gin.Engine  // the http handler
	health  bool         // server state for health checks
	started time.Time    // when the server started
}

var sensors = map[string]Sensor{
	"1": {Name: "1", Location: Coordinates{}, Tags: SensorTags{}},
	"2": {Name: "2", Location: Coordinates{}, Tags: SensorTags{}},
	"3": {Name: "3", Location: Coordinates{}, Tags: SensorTags{}},
}

func New() (server *Server) {
	ginRouter := gin.Default()
	server = &Server{
		srv:    &http.Server{Handler: ginRouter},
		router: ginRouter,
		health: false,
	}

	server.setupRoutes()
	return server
}

// TODO
func (s *Server) Serve() (err error) {
	s.SetStatus(true)
	s.started = time.Now()

	log.Println("server starting")
	if err = s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// TODO
func (s *Server) Shutdown() {
	log.Println("gracefully shutting down server")
	s.SetStatus(false)
	s.srv.Shutdown(context.Background())
	log.Println("server successfully shutdown")
}

// TODO
func (s *Server) SetStatus(health bool) {
	s.Lock()
	s.health = health
	s.Unlock()
}

// TODO
func (s *Server) setupRoutes() {
	s.router.GET("/allsensors", s.ListSensors)
	s.router.GET("/sensor:name", s.GetSensor)
	s.router.POST("/sensor", s.InsertSensor)
	s.router.PUT("/sensor", s.UpdateSensor)
	s.router.GET("/nearest", s.NearestLocation)
	s.router.GET("/health", s.Health)
}
