package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	HTTPDconfig httpdConfig
	serveHTTP   = flag.Bool("serve", false, "Run as webserver / cloud function")
)

func main() {
	if os.Getenv("NC_SERVE_HTTP") != "" {
		*serveHTTP = true
	}
	flag.Parse()
	HTTPDconfig.address = thisOrThat(os.Getenv("NC_TCP_PORT"), ":2001")
	fmt.Printf("NC booting. Webeserver: %v\n", *serveHTTP)
	RunWebserver(HTTPDconfig)
}
