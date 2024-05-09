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
			Name:   "",
			Usage:  "",
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
