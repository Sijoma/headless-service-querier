package pinger

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type pinger struct {
	http.Client
}

func newPinger() *pinger {
	return &pinger{
		http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (f *pinger) do(host string, port string) error {
	resp, err := f.Get(fmt.Sprintf("http://%v:"+port, host))
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error while reading response: %v\n", err)
	}
	if raw != nil {
		fmt.Printf("Output from ip: %v, %v \n\n", host, string(raw))
	}
	return nil
}

func New(hostName, port string) error {
	if hostName == "" {
		return fmt.Errorf("hostname is empty")
	}
	pingerInstance := newPinger()
	for {
		log.Printf("connecting to %s", hostName)
		time.Sleep(time.Second * 2)
		ips, err := net.LookupIP(hostName)
		if err != nil {
			fmt.Printf("error while looking up ips from headless service: %v\n", err)
			continue
		}
		log.Printf("found %d ips, %v", len(ips), ips)
		for _, ip := range ips {
			fmt.Printf("Calling: %v \n", ip.String())
			time.Sleep(2 * time.Second)
			err := pingerInstance.do(ip.String(), port)
			if err != nil {
				fmt.Printf("error while pinging endpoint %v: %v\n", ip.String(), err)
			}
		}
	}

}
