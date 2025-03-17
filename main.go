package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"headless/pinger"
	"headless/server"
)

func main() {
	serverMode := flag.Bool("server", true, "run in server mode")
	hostName := flag.String("hostname", "localhost", "hostname")
	port := flag.Int("port", 8080, "port")
	flag.Parse()
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		podName = "pod"
	}

	fmt.Printf("Server mode: %v\n", *serverMode)
	fmt.Printf("Hostname: %v\n", *hostName)
	fmt.Printf("Port: %v\n", *port)
	fmt.Printf("Pod name: %v\n", podName)

	if *serverMode {
		server.New(podName, *port, *hostName)
	} else {
		err := pinger.New(podName, *hostName, strconv.Itoa(*port))
		if err != nil {
			log.Fatal(err)
		}
	}
}
