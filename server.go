package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var version = "0.0.1"

type Options struct {
	// Port Number (default=8080)
	Port int

	// Path of custom 404.html file [only file path]
	Error404 string

	// Slice of all the paths to ignore, if any request is made to access the ignored file(s), it will be redirected to 404.html [only file path(s)]
	IgnorePath []string
}

var Opts Options

func errorView(w http.ResponseWriter, r *http.Request) {
	htmlTmp := ""
	if Opts.Error404 != "" {
		rf, err := ioutil.ReadFile(Opts.Error404)
		if err != nil {
			fmt.Println(err)
		} else {
			htmlTmp = string(rf)
		}
	} else {
		htmlTmp = fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta http-equiv="X-UA-Compatible" content="IE=edge">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Error -Page Not Found</title>
		</head>
		<style>
			body {
				font-family: 'Lucida Console';
				background-color: rgb(248, 248, 248);
				font-weight: 100;
			}
		
			a {
				color: rgb(38, 77, 206);
				text-decoration: none;
				cursor: pointer;
			}
		
			#content {
				margin: auto;
				text-align: center;
			}
		
			h1 {
				font-size: 5rem;
			}
		</style>
		<body>
			<div id="content">
				<h1>404</h1>
				<h3>It looks like we couldn't find the page</h3>
				<h3>Try going <a href="%v">Home</a> or <a onclick="window.history.back()">Back</a></h3>
			</div>
		</body>
		</html>`, "/")
	}
	
	fmt.Fprint(w, htmlTmp)
}

func fileServerC(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			errorView(rw, r)
		}

		for _, p := range Opts.IgnorePath {
			if p == r.URL.Path {
				errorView(rw, r)
			}
		}
		fsh.ServeHTTP(rw, r)
	})
}

func Start() {
	fmt.Printf("Server is started at port:%v\n", Opts.Port)
	fmt.Printf("Open 127.0.0.1:%v\n", Opts.Port)
	fmt.Println("Press Ctrl+C to terminate")
	
	fp, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	http.Handle("/", fileServerC(http.Dir(fp)))
	serv := fmt.Sprintf(":%v", Opts.Port)
    fmt.Print(http.ListenAndServe(serv, nil))
}

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
					Opts = Options{
						Port: c.Int("port"),
						Error404: c.String("404"),
						IgnorePath: c.StringSlice("ignore"),
					}
					Start()
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