package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
)

type worker struct {
	name string
}

func New(name string, port int) {
	http.HandleFunc("/", worker{name: name}.workerCountHandler)
	log.Printf("Server starting on port %v \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func (m worker) workerCountHandler(w http.ResponseWriter, r *http.Request) {
	workers := rand.Intn(5)
	w.Header().Set("Content-Type", "application/json")
	slog.Info("serving", slog.String("endpoint", r.URL.String()), slog.Int("workerCount", workers))
	response, _ := json.Marshal(workersResp{WorkersCount: workers, WorkerName: m.name})
	w.Write(response)
}

type workersResp struct {
	WorkersCount int    `json:"workersCount"`
	WorkerName   string `json:"workerName"`
}
