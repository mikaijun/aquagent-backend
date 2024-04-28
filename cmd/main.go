package main

import (
	"flag"

	"github.com/mikaijun/anli/pkg/interfaces/api/server"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()
}

func main() {
	server.Serve(addr)
}
