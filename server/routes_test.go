package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
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
	suite.srv = *server.New("fakeaddress")
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
	suite.testContext.Request = &http.Request{
		Header: make(http.Header),
	}
	suite.testContext.Request.Method = "POST"
	suite.testContext.Request.Header.Set("Content-Type", "application/json")

	sensor := &server.Sensor{
		Name: "foo",
		Location: server.Coordinates{
			Latitude:  0,
			Longitude: 0,
		},
		Tags: server.SensorTags{
			Name: "foo",
			Unit: "bar",
		},
	}
	bodyBytes, err := json.Marshal(sensor)
	if err != nil {
		panic(err)
	}
	body := bytes.NewBuffer(bodyBytes)
	suite.testContext.Request.Body = io.NopCloser(body)

	suite.srv.AddSensor(suite.testContext)
	suite.Equal(201, suite.responseRecorder.Code)
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
	suite.testContext.Params = []gin.Param{
		{
			Key:   "lat",
			Value: "30",
		},
		{
			Key:   "lon",
			Value: "100",
		},
	}
	suite.srv.NearestSensor(suite.testContext)
	suite.Equal(200, suite.responseRecorder.Code)

	body, err := io.ReadAll(suite.responseRecorder.Body)
	suite.Nil(err)

	var responseSensor server.Sensor
	suite.Nil(json.Unmarshal(body, &responseSensor))
	suite.Equal(responseSensor.Name, "L1ANG")

	suite.testContext.Params = []gin.Param{
		{
			Key:   "lat",
			Value: "40",
		},
		{
			Key:   "lon",
			Value: "170",
		},
	}
	suite.srv.NearestSensor(suite.testContext)
	suite.Equal(200, suite.responseRecorder.Code)

	body, err = io.ReadAll(suite.responseRecorder.Body)
	suite.Nil(err)

	suite.Nil(json.Unmarshal(body, &responseSensor))
	suite.Equal(responseSensor.Name, "C1MAG")
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
