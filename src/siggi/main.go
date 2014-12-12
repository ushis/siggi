package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sighub"
	"syscall"
)

var (
	root    string
	address string
)

func init() {
	flag.StringVar(&root, "root", ".", "root directory")
	flag.StringVar(&address, "listen", ":8080", "address to listen to")
}

func main() {
	flag.Parse()

	listener, err := listen()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Fprintf(os.Stderr, "Listening on %s\n", address)

	hub := sighub.New()
	go hub.Run()
	defer hub.Die()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(root)))
	mux.Handle("/socket", hub.HTTPHandler())

	go http.Serve(listener, mux)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
}

func listen() (net.Listener, error) {
	if len(address) > 0 && address[0] == '/' {
		return net.Listen("unix", address)
	}
	return net.Listen("tcp", address)
}
