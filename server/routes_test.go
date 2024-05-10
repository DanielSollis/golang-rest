package server_test

import (
	"io"
	"math"
	"net/http/httptest"
	"pingthings/server"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	srv              server.Server
	responseRecorder *httptest.ResponseRecorder
	testContext      *gin.Context
}

func (suite *testSuite) SetupTest() {
	suite.setupRecorder()
	suite.srv = *server.New("foo")
}

func (suite *testSuite) setupRecorder() {
	recorder := httptest.NewRecorder()
	suite.responseRecorder = recorder
	suite.testContext, _ = gin.CreateTestContext(recorder)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestListSensors() {
	suite.srv.ListSensors(suite.testContext)
	suite.Equal(200, suite.responseRecorder.Code)

	body, err := io.ReadAll(suite.responseRecorder.Body)
	suite.Nil(err)
	suite.NotNil(body)
}

func (suite *testSuite) TestAddSensor() {
	// TODO
}

func (suite *testSuite) TestGetSensor() {
	suite.testContext.Params = []gin.Param{{
		Key:   "name",
		Value: "L1MAG",
	}}
	suite.srv.GetSensor(suite.testContext)
	suite.Equal(200, suite.responseRecorder.Code)

	suite.setupRecorder()
	suite.testContext.Params = []gin.Param{{
		Key:   "name",
		Value: "NOTASENSOR",
	}}
	suite.srv.GetSensor(suite.testContext)
	suite.Equal(404, suite.responseRecorder.Code)
}

func (suite *testSuite) TestNearestSensor() {
	// TODO
}

func (suite *testSuite) TestStatusCheck() {
	suite.srv.StatusCheck(suite.testContext)
	suite.Equal(200, suite.responseRecorder.Code)
}

func (suite *testSuite) TestHaversine() {
	// test one
	userCoordinates := server.Coordinates{
		Latitude:  0,
		Longitude: 0,
	}
	sensorCoordinates := server.Coordinates{
		Latitude:  0,
		Longitude: 180,
	}
	distance := server.Haversine(userCoordinates, sensorCoordinates)
	expected := 20015
	suite.EqualValues(expected, math.Round(distance))

	// Test two
	userCoordinates = server.Coordinates{
		Latitude:  51.5007,
		Longitude: 0.1246,
	}
	sensorCoordinates = server.Coordinates{
		Latitude:  40.6892,
		Longitude: 74.0445,
	}
	distance = server.Haversine(userCoordinates, sensorCoordinates)
	expected = 5575
	suite.EqualValues(expected, math.Round(distance))
}
