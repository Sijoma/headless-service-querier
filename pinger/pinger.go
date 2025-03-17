package pinger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"headless/worker"
)

type Pinger struct {
	responder      string
	hostName, port string
	client         http.Client
}

func New(responder, hostName string, port string) *Pinger {
	return &Pinger{
		responder: responder,
		hostName:  hostName,
		port:      port,
		client:    http.Client{Timeout: time.Second * 5},
	}
}

func (f *Pinger) do(hostName string) (*BackendResponse, error) {
	resp, err := f.client.Get(fmt.Sprintf("http://%v:%v", hostName, f.port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response worker.CountResponse
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &BackendResponse{
		Name: response.WorkerName,
		Info: fmt.Sprintf("%v workers", response.WorkersCount),
	}, nil
}

func (f *Pinger) Ping() ([]*BackendResponse, error) {
	log.Printf("connecting to %s", f.hostName)
	ips, err := net.LookupIP(f.hostName)
	if err != nil {
		return nil, fmt.Errorf("error while looking up ips from headless service: %w", err)
	}
	log.Printf("found %d ips, %v", len(ips), ips)
	var result []*BackendResponse
	for _, ip := range ips {
		fmt.Printf("Calling: %v \n", ip.String())

		resp, err := f.do(ip.String())
		if err != nil {
			return nil, fmt.Errorf("error while pinging %v: %w", ip, err)
		}
		result = append(result, resp)
	}
	return result, nil
}

func (f *Pinger) Handler(w http.ResponseWriter, _ *http.Request) {
	resp, err := f.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(response{
		Responder:   f.responder,
		BackendInfo: resp,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		return
	}
	return
}

type response struct {
	Responder   string             `json:"responder"`
	BackendInfo []*BackendResponse `json:"backend_info"`
}

type BackendResponse struct {
	Name string `json:"pod_name"`
	Info string `json:"info"`
}
