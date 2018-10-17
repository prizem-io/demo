package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cenkalti/backoff"
)

const iptablesFormat = `iptables -t nat -%s PREROUTING -s %s/16 -p tcp --dport 80 -j DNAT --to %s:%s -m comment --comment "prizem/dnat-to-proxy"`

func main() {
	proxyHost := os.Getenv("PROXY_HOST")
	if proxyHost == "" {
		proxyHost = "proxy"
	}
	proxyPort := os.Getenv("PROXY_PORT")
	if proxyPort == "" {
		proxyPort = "50062"
	}

	proxyIP := os.Getenv("PROXY_IP")

	notify := func(err error, d time.Duration) {
		fmt.Printf("Failed attempt: %v -> will retry in %s\n", err, d)
	}

	var addrs []string
	if proxyIP != "" {
		addrs = []string{proxyIP}
	} else {
		err := backoff.RetryNotify(func() (err error) {
			// TODO: Figure out how to inject the IP of the proxy with `network_mode: host`
			addrs, err = net.LookupHost(proxyHost)
			return
		}, backoff.NewExponentialBackOff(), notify)
		if err != nil {
			panic(err)
		}
	}

	if len(addrs) == 0 {
		panic(fmt.Errorf("no addresses were returned for %q", proxyHost))
	}

	ip := net.IPv4zero
	fmt.Printf("IP Addresses for %q:\n", proxyHost)
	for _, addr := range addrs {
		checkIP := net.ParseIP(addr)
		fmt.Printf("- %s\n", addr)
		if p4 := checkIP.To4(); len(p4) == net.IPv4len {
			ip = checkIP
		}
	}
	fmt.Println()

	if ip.IsUnspecified() {
		panic(fmt.Errorf("no IPv4 addresses were returned for %q", proxyHost))
	}

	cidr := ip.Mask(net.IPv4Mask(0xff, 0xff, 0, 0))
	cmdString := fmt.Sprintf(iptablesFormat, "A", cidr, ip, proxyPort)
	fmt.Println(cmdString)
	cmdParts := strings.Split(cmdString, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	cmdString = fmt.Sprintf(iptablesFormat, "D", cidr, ip, proxyPort)
	fmt.Println(cmdString)
	cmdParts = strings.Split(cmdString, " ")
	cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
