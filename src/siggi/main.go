package main

import (
	"flag"
	"log"
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
		log.Fatalln(err)
	}
	defer listener.Close()

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
