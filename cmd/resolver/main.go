package main

import (
	"context"
	"os"

	"github.com/1inch/p2p-network/relayer"
	"github.com/urfave/cli"
)

// TODO: setup cli interface
func main() {
	app := &cli.App{
		Name:  "resolver",
		Usage: "Resolver node",
		Commands: []cli.Command{
			{
				Name:  "run",
				Usage: "Runs resolver node",
				Action: func(c *cli.Context) error {
					node, err := relayer.New()
					if err != nil {
						// TODO: handle error
					}

					node.Run(context.Background())

					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
