package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv     *http.Server // http server for API defaults
	router  *gin.Engine  // the http handler
	healthy bool         // server state for health checks
	started time.Time    // when the server started
}

var sensors = map[string]Sensor{
	"L1MAG": {
		Name: "L1MAG",
		// Null Island off the coast of Africa
		Location: Coordinates{
			Latitude:  0,
			Longitude: 0,
		},
		Tags: SensorTags{
			Name:      "L1MAG",
			Unit:      "volts",
			Ingress:   "",
			Distiller: "",
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
			Name:      "L1ANG",
			Unit:      "deg",
			Ingress:   "",
			Distiller: "",
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
			Name:      "C1MAG",
			Unit:      "amps",
			Ingress:   "",
			Distiller: "",
		},
	},
}

func New(addr string) (server *Server) {
	ginRouter := gin.Default()
	server = &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: ginRouter,
		},
		router:  ginRouter,
		healthy: false,
	}

	server.setupRoutes()
	return server
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
	s.router.GET("/allsensors", s.ListSensors)
	s.router.POST("/sensor", s.AddSensor)
	s.router.GET("/sensor/:name", s.GetSensor)
	s.router.GET("/nearest/:lat/:lon", s.NearestSensor)
	s.router.GET("/health", s.StatusCheck)
}
