package main

import (
	"fmt"
	"os"

	"github.com/soubikbhuiwk007/gnet/internal/serve"
	"github.com/urfave/cli/v2"
)

var version = "0.0.1"

func main() {
	(&cli.App{
		Name: "gnet",
		Version: version,
		Description: "Simple Web Server in Go",
		Usage: "Simple Web Server in Go",
		UsageText: "gnet [command] [opts...] [args...]",
		CustomAppHelpTemplate: fmt.Sprintf(`Simple Web Server in Go

VERSION: %v

COMMANDS:
		start	Start the server
		help, h Shows a list of commands or help for one command

SERVER OPTIONS:
		--port		Give the port number (default=8080)
		--404		Set a custom 404 error page
		--ignore		Set a list of files/folders to ignore, this will load 404.html if any request is made

GLOBAL OPTIONS:
		--help, -h	show help
		--version, -v	print the version
`, version),
		Commands: []*cli.Command{
			{
				Name: "start",
				Description: "Start a server",
				Usage: "Start a server",
				UsageText: "start [opts...] [args...]",
				Action: func(c *cli.Context) error {
					serve.Opts = serve.Options{
						Port: c.Int("port"),
						Error404: c.String("404"),
						IgnorePath: c.StringSlice("ignore"),
					}
					serve.Start()
					return nil
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name: "port",
						Usage: "Port Number",
						Value: 8080,
						Required: false,
					},
					&cli.StringFlag{
						Name: "404",
						Usage: "Set custom 404.html",
						Value: "",
						Required: false,
					},
					&cli.StringSliceFlag{
						Name: "ignore",
						Usage: "Set the folder and file paths to ignore. This will load 404.html if any request if made",
						Required: false,
					},
				},
			},
		},
	}).Run(os.Args)
}