package serve

type Options struct {
	// Port Number (default=8080)
	Port int

	// Path of custom 404.html file [only file path]
	Error404 string

	// Slice of all the paths to ignore, if any request is made to access the ignored file(s), it will be redirected to 404.html [only file path(s)]
	IgnorePath []string
}