package serve

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

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