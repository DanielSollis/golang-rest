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
	srv     *http.Server // http server for API defaults
	router  *gin.Engine  // http handler
	db      *store       // SQLite connection
	healthy bool         // server state for health checks
	started time.Time    // when the server started
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
	Unit string `json:"unit"`
	Name string `json:"name"`
}

func New(addr string) (server *Server, err error) {
	ginRouter := gin.Default()
	server = &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: ginRouter,
		},
		router:  ginRouter,
		healthy: false,
	}

	if server.db, err = initDB(); err != nil {
		return nil, err
	}

	server.setupRoutes()

	return server, nil
}

func (s *Server) Serve() (err error) {
	s.healthy = true
	s.started = time.Now()

	// Listen for sigint call to close server
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

func (s *Server) shutdown() {
	s.db.conn.Close()
	ctx := context.Background()
	_ = s.srv.Shutdown(ctx)
}

func (s *Server) setupRoutes() {
	s.router.GET("/allsensors", s.listSensors)
	s.router.POST("/sensor", s.addSensor)
	s.router.POST("/sensor/:name", s.updateSensor)
	s.router.GET("/sensor/:name", s.getSensor)
	s.router.GET("/nearest/:lat/:lon", s.getNearestSensor)
	s.router.GET("/health", s.statusCheck)
}
