package client

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "",
			Aliases: []string{},
			Usage:   "",
			Value:   "",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Action: serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "",
					Aliases: []string{},
					Usage:   "",
				},
			},
		},
		{
			Name:   "shutdown",
			Action: serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "",
					Aliases: []string{},
					Usage:   "",
				},
			},
		},
		{
			Name:   "getAll",
			Action: serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "",
					Aliases: []string{},
					Usage:   "",
				},
			},
		},
		{
			Name:   "insert",
			Action: serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "",
					Aliases: []string{},
					Usage:   "",
				},
			},
		},
	}

	app.Run(os.Args)
}

func serve(c *cli.Context) (err error) {
	return nil
}

func shutdown(c *cli.Context) (err error) {
	return nil
}

func insert(c *cli.Context) (err error) {
	return nil
}

func getAll(c *cli.Context) (err error) {
	return nil
}
