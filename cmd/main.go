package main

import (
	"flag"

	server "github.com/mikaijun/aquagent/pkg/interfaces"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8000", "tcp host:port to connect")
	flag.Parse()
}

func main() {
	server.Serve(addr)
}
