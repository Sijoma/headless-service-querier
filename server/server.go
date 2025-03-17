package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"headless/pinger"
	"headless/worker"
)

func New(name string, port int, headlessServiceHost string) {
	w := worker.NewWorker(name)
	http.HandleFunc("/", w.CountHandler)

	pingInstance := pinger.New(name, headlessServiceHost, strconv.Itoa(port))

	// Api/Workers allows to get runtime information of the workers
	http.HandleFunc("/api/workers", pingInstance.Handler)

	log.Printf("Server starting on port %v \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
