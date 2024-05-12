package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv     *http.Server // http server for API defaults.
	gin     *gin.Engine  // http handler.
	db      *store       // SQLite connection.
	healthy bool         // server state for health checks.
	started time.Time    // when the server started.
}

type Sensor struct {
	Name     string      `json:"name"`
	Location Coordinates `json:"location"`
	Tags     SensorTags  `json:"tags"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SensorTags struct {
	Name      string `json:"name"`
	Unit      string `json:"unit"`
	Ingress   string `json:"ingress"`
	Distiller string `json:"distiller"`
}

func New(addr string) (server *Server, err error) {
	ginEngine := gin.Default()
	server = &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: ginEngine,
		},
		gin:     ginEngine,
		healthy: false,
	}

	if server.db, err = newStore(); err != nil {
		return nil, err
	}

	server.setupRoutes()

	return server, nil
}

func (s *Server) Serve() (err error) {
	s.healthy = true
	s.started = time.Now()

	// Listen for sigint call to gracefully
	// shutdown the server.
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGINT)
	go func() {
		<-interupt
		log.Println("shutting down server")
		s.shutdown()
	}()

	log.Println("server starting")
	if err = s.srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Gracefully shutdown the server. Closes
// the database connection and http server.
func (s *Server) shutdown() {
	s.db.conn.Close()
	ctx := context.Background()
	_ = s.srv.Shutdown(ctx)
}

// Add REST endpoints to the Gin engine.
func (s *Server) setupRoutes() {
	s.gin.GET("/allsensors", s.listSensors)
	s.gin.GET("/sensor/:name", s.getSensor)
	s.gin.POST("/sensor", s.addSensor)
	s.gin.PUT("/sensor/:name", s.updateSensor)
	s.gin.GET("/nearest/:lat/:lon", s.getNearestSensor)
	s.gin.GET("/health", s.statusCheck)
}
