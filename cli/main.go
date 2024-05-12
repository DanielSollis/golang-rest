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
					Aliases:  []string{"la"},
					Required: true,
				},
				&cli.Float64Flag{
					Name:     "lon",
					Aliases:  []string{"lo"},
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
			Name:     "update",
			Category: "client",
			Action:   updateSensor,
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
					Aliases:  []string{"la"},
					Required: true,
				},
				&cli.Float64Flag{
					Name:     "lon",
					Aliases:  []string{"lo"},
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
					Name:    "endpoint",
					Aliases: []string{"e"},
					Value:   "http://localhost:8080/nearest",
				},
				&cli.Float64Flag{
					Name:     "lat",
					Aliases:  []string{"la"},
					Required: true,
				},
				&cli.Float64Flag{
					Name:     "lon",
					Aliases:  []string{"lo"},
					Required: true,
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
	if srv, err = server.New(addr); err != nil {
		return err
	}

	if err := srv.Serve(); err != nil {
		return err
	}
	return nil
}

func listSensors(c *cli.Context) (err error) {
	url := c.String("endpoint")
	var responseString string
	if responseString, err = getEndpoint(url); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func addSensor(c *cli.Context) (err error) {
	name, unit := c.String("name"), c.String("unit")
	lat, lon := c.Float64("lat"), c.Float64("lon")
	sensor := server.CreateSensor(name, unit, lat, lon)

	var responseString string
	url := c.String("endpoint")
	if responseString, err = postEndpoint(url, sensor); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func updateSensor(c *cli.Context) (err error) {
	name, unit := c.String("name"), c.String("unit")
	lat, lon := c.Float64("lat"), c.Float64("lon")
	sensor := server.CreateSensor(name, unit, lat, lon)

	var responseString string
	url := fmt.Sprintf("%s/%s", c.String("endpoint"), name)
	if responseString, err = postEndpoint(url, sensor); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func getSensor(c *cli.Context) (err error) {
	var responseString string
	endpoint, name := c.String("endpoint"), c.String("name")
	url := fmt.Sprintf("%s/%s", endpoint, name)
	if responseString, err = getEndpoint(url); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func nearestSensor(c *cli.Context) (err error) {
	var responseString string
	lat, lon := c.Float64("lat"), c.Float64("lon")
	url := fmt.Sprintf("%s/%f/%f", c.String("endpoint"), lat, lon)
	if responseString, err = getEndpoint(url); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func statusCheck(c *cli.Context) (err error) {
	endpoint := c.String("endpoint")
	var responseString string
	if responseString, err = getEndpoint(endpoint); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseString)
	return nil
}

func getEndpoint(url string) (_ string, err error) {
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

func postEndpoint(url string, toMarshal interface{}) (_ string, err error) {
	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(toMarshal); err != nil {
		return "", err
	}

	var response *http.Response
	requestBody := bytes.NewBuffer(jsonBytes)
	if response, err = http.Post(url, "application/json", requestBody); err != nil {
		return "", err
	}

	var responseBody []byte
	if responseBody, err = io.ReadAll(response.Body); err != nil {
		return "", err
	}
	return string(responseBody), nil
}
