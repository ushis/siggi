package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	listenAddr string
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8080", "address to listen to")
}

func main() {
	flag.Parse()

	l, err := listen(listenAddr)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer l.Close()

	go http.Serve(l, NewHub().HTTPHandler())

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sig
}

func listen(addr string) (net.Listener, error) {
	if len(addr) > 0 && addr[0] == '/' {
		return net.Listen("unix", addr)
	}
	return net.Listen("tcp", addr)
}
