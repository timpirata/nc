package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	serveHTTP = flag.Bool("serve", false, "Run as webserver / cloud function")
)

func main() {
	if os.Getenv("NC_SERVE_HTTP") != "" {
		*serveHTTP = true
	}
	flag.Parse()

	httpdConfig := httpdConfig{
		address: thisOrThat(os.Getenv("NC_TCP_PORT"), ":2001"),
		enabled: *serveHTTP,
	}
	fmt.Println("NC booting...")
	RunWebserver(httpdConfig)
}
