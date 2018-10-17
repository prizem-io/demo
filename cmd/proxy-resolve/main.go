package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
)

func main() {
	proxyFile := os.Getenv("PROXY_FILE")
	if proxyFile == "" {
		proxyFile = "etc/proxy-ip.txt"
	}
	proxyHost := os.Getenv("PROXY_HOST")
	if proxyHost == "" {
		proxyHost = "proxy"
	}

	notify := func(err error, d time.Duration) {
		fmt.Printf("Failed attempt: %v -> will retry in %s\n", err, d)
	}

	var addrs []string
	err := backoff.RetryNotify(func() (err error) {
		addrs, err = net.LookupHost(proxyHost)
		return nil
	}, backoff.NewExponentialBackOff(), notify)
	if err != nil {
		panic(err)
	}

	ipv4Addrs := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		checkIP := net.ParseIP(addr)
		if p4 := checkIP.To4(); len(p4) == net.IPv4len {
			ipv4Addrs = append(ipv4Addrs, addr)
		}
	}

	fmt.Printf("Proxy IPs are %s\n", strings.Join(ipv4Addrs, ", "))

	err = ioutil.WriteFile(proxyFile, []byte(strings.Join(ipv4Addrs, ",")), 0777)
	if err != nil {
		panic(err)
	}
}
