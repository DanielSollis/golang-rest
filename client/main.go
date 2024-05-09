package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pingthings/server"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Category: "server",
			Action:   serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Aliases: []string{"a"},
					Value:   "localhost:8080",
				},
			},
		},
		{
			Name:     "list",
			Category: "client",
			Action:   listSensors,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "endpoint",
					Aliases: []string{"e"},
					Value:   "http://localhost:8080/allsensors",
				},
			},
		},

		{
			Name:     "add",
			Category: "client",
			Action:   addSensor,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "endpoint",
					Aliases: []string{"e"},
					Value:   "http://localhost:8080/sensor",
				},
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Required: true,
				},
				&cli.Float64Flag{
					Name:     "lat",
					Aliases:  []string{"a"},
					Required: true,
				},
				&cli.Float64Flag{
					Name:     "lon",
					Aliases:  []string{"o"},
					Required: true,
				},
				&cli.StringFlag{
					Name:     "unit",
					Aliases:  []string{"u"},
					Required: true,
				},
			},
		},
		{
			Name:     "get",
			Category: "client",
			Action:   getSensor,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "endpoint",
					Aliases: []string{"e"},
					Value:   "http://localhost:8080/sensor",
				},
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Required: true,
				},
			},
		},
		{
			Name:     "nearest",
			Category: "client",
			Action:   nearestSensor,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "",
					Usage: "",
				},
			},
		},
		{
			Name:     "status",
			Category: "client",
			Action:   statusCheck,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "endpoint",
					Aliases: []string{"e"},
					Value:   "http://localhost:8080/health",
				},
			},
		},
	}

	app.Run(os.Args)
}

var srv *server.Server

func serve(c *cli.Context) (err error) {
	addr := c.String("address")
	srv = server.New(addr)
	if err := srv.Serve(); err != nil {
		fmt.Println(err)
	}
	return nil
}

func listSensors(c *cli.Context) (err error) {
	endpoint := c.String("endpoint")
	var responseString string
	if responseString, err = endpointGet(endpoint); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func addSensor(c *cli.Context) (err error) {
	sensor := server.Sensor{
		Name: c.String("name"),
		Location: server.Coordinates{
			Latitude:  c.Float64("lat"),
			Longitude: c.Float64("lat"),
		},
		Tags: server.SensorTags{
			Name: c.String("name"),
			Unit: c.String("unit"),
		},
	}

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(sensor); err != nil {
		fmt.Println(err)
		return err
	}

	url := c.String("endpoint")
	requestBody := bytes.NewBuffer(jsonBytes)
	var response *http.Response
	if response, err = http.Post(url, "application/json", requestBody); err != nil {
		fmt.Println(err)
		return err
	}

	var responseBody []byte
	if responseBody, err = io.ReadAll(response.Body); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(responseBody))
	return nil
}

func getSensor(c *cli.Context) (err error) {
	url := fmt.Sprintf("%s/%s", c.String("endpoint"), c.String("name"))
	var responseString string
	if responseString, err = endpointGet(url); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func nearestSensor(c *cli.Context) (err error) {
	return nil
}

func statusCheck(c *cli.Context) (err error) {
	endpoint := c.String("endpoint")
	var responseString string
	if responseString, err = endpointGet(endpoint); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func endpointGet(url string) (_ string, err error) {
	var response *http.Response
	if response, err = http.Get(url); err != nil {
		return "", err
	}

	var body []byte
	if body, err = io.ReadAll(response.Body); err != nil {
		return "", err
	}
	return string(body), nil
}
