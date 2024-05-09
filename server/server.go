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
	"1": {
		Name: "1",
		// Null Island off the coast of Africa
		Location: Coordinates{
			Latitude:  0,
			Longitude: 0,
		},
		Tags: SensorTags{
			Name: "foo",
			Unit: "foo",
		},
	},
	"2": {
		Name: "2",
		Location: Coordinates{
			Latitude:  0,
			Longitude: 0,
		},
		Tags: SensorTags{
			Name: "foo",
			Unit: "foo",
		},
	},
	"3": {
		Name: "3",
		Location: Coordinates{
			Latitude:  0,
			Longitude: 0,
		},
		Tags: SensorTags{
			Name: "foo",
			Unit: "foo",
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

// TODO
func (s *Server) Serve() (err error) {
	s.healthy = true
	s.started = time.Now()

	log.Println("server starting")
	if err = s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// TODO
func (s *Server) setupRoutes() {
	s.router.GET("/allsensors", s.ListSensors)
	s.router.GET("/sensor/:name", s.GetSensor)
	s.router.POST("/sensor", s.InsertSensor)
	s.router.PUT("/sensor", s.UpdateSensor)
	s.router.GET("/nearest", s.NearestLocation)
	s.router.GET("/health", s.Health)
}
