package main

import (
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
					Name:    "",
					Aliases: []string{},
					Usage:   "",
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
	var response *http.Response
	if response, err = http.Get(endpoint); err != nil {
		fmt.Println(err)
		return err
	}

	var body []byte
	if body, err = io.ReadAll(response.Body); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}

func addSensor(c *cli.Context) (err error) {
	return nil
}

func getSensor(c *cli.Context) (err error) {
	var response *http.Response
	url := fmt.Sprintf("%s/%s", c.String("endpoint"), c.String("name"))
	if response, err = http.Get(url); err != nil {
		fmt.Println(err)
		return err
	}

	var body []byte
	if body, err = io.ReadAll(response.Body); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}

func nearestSensor(c *cli.Context) (err error) {
	return nil
}

func statusCheck(c *cli.Context) (err error) {
	return nil
}
