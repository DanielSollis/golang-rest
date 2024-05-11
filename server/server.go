package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv     *http.Server // http server for API defaults
	router  *gin.Engine  // the http handler
	db      *sql.DB      // SQLite connection
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

var sensors = map[string]*Sensor{
	"L1MAG": {
		Name: "L1MAG",
		// Null Island off the coast of Africa
		Location: Coordinates{
			Latitude:  0,
			Longitude: 0,
		},
		Tags: SensorTags{
			Name: "L1MAG",
			Unit: "volts",
		},
	},
	"L1ANG": {
		Name: "L1ANG",
		// Disneyland, Anaheim
		Location: Coordinates{
			Latitude:  33.8,
			Longitude: 117.9,
		},
		Tags: SensorTags{
			Name: "L1ANG",
			Unit: "deg",
		},
	},
	"C1MAG": {
		Name: "C1MAG",
		// Hobbiton, New Zealand
		Location: Coordinates{
			Latitude:  37.8,
			Longitude: 175.7,
		},
		Tags: SensorTags{
			Name: "C1MAG",
			Unit: "amps",
		},
	},
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

	log.Println("server starting")
	if err = s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) setupRoutes() {
	s.router.GET("/allsensors", s.listSensors)
	s.router.POST("/sensor", s.addSensor)
	s.router.GET("/sensor/:name", s.getSensor)
	s.router.GET("/nearest/:lat/:lon", s.nearestSensor)
	s.router.GET("/health", s.statusCheck)
}
