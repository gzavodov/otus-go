package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gzavodov/otus-go/telnet"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var host string
	var port string
	timeout := flag.Duration("timeout", time.Duration(10)*time.Second, "Connection Timeout (time duration: 10s, 1m etc.)")

	flag.Parse()
	args := flag.Args()

	if len(args) > 1 {
		host = args[0]
		port = args[1]
	}

	if host == "" || port == "" {
		fmt.Println("Usage:", os.Args[0], "--timeout=[timeout] [host] [port]")
		fmt.Println("Help:", os.Args[0], "-help")
		return
	}

	client := telnet.NewClient(
		"tcp",
		net.JoinHostPort(host, port),
		*timeout,
		os.Stdin,
		os.Stdout,
	)

	err := client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
