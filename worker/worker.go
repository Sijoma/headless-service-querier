package worker

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
)

type Worker struct {
	name string
}

func (m Worker) CountHandler(w http.ResponseWriter, r *http.Request) {
	workers := rand.Intn(5)
	w.Header().Set("Content-Type", "application/json")
	slog.Info("serving", slog.String("endpoint", r.URL.String()), slog.Int("workerCount", workers))
	response, _ := json.Marshal(CountResponse{WorkersCount: workers, WorkerName: m.name})
	w.Write(response)
}

type CountResponse struct {
	WorkersCount int    `json:"workersCount"`
	WorkerName   string `json:"workerName"`
}

func NewWorker(name string) Worker {
	return Worker{name: name}
}
