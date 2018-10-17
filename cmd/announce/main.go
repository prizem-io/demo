package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/cenkalti/backoff"
	"github.com/satori/go.uuid"
)

const registerTemplate = `{"id": "%s", "service": "%s", "name": "%s", "ports": [{"port": %d, "protocol": "HTTP/1"}]}'`

func main() {
	registerURI := os.Getenv("REGISTER_URI")
	if registerURI == "" {
		registerURI = "http://proxy:6060/register"
	}

	servicesString := os.Getenv("SERVICES")
	services := strings.FieldsFunc(servicesString, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})

	notify := func(err error, d time.Duration) {
		fmt.Printf("Failed attempt: %v -> will retry in %s\n", err, d)
	}
	ebo := backoff.NewExponentialBackOff()
	client := &http.Client{}

	for _, service := range services {
		ebo.Reset()
		var ips []net.IP
		err := backoff.RetryNotify(func() (err error) {
			ips, err = getHostIPv4s(service)
			return
		}, ebo, notify)
		if err != nil {
			panic(err)
		}

		payload := fmt.Sprintf(registerTemplate, uuid.NewV4(), service, service, 8000)
		req, err := http.NewRequest("POST", registerURI, bytes.NewReader([]byte(payload)))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		for _, ip := range ips {
			req.Header.Set("X-Target-Host", ip.String())

			ebo.Reset()
			err = backoff.RetryNotify(func() (err error) {
				// Do the request
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					panic(fmt.Errorf("register return status %d", resp.StatusCode))
				}
				return nil
			}, ebo, notify)
			if err != nil {
				panic(err)
			}
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	req, err := http.NewRequest("DELETE", registerURI, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func getHostIPv4s(host string) ([]net.IP, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, err
	}

	results := make([]net.IP, 0, len(addrs))

	for _, addr := range addrs {
		checkIP := net.ParseIP(addr)
		if p4 := checkIP.To4(); len(p4) == net.IPv4len {
			results = append(results, checkIP)
		}
	}

	return results, nil
}
