package main

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
)

func validateVaultHost(host string) bool {
	var hostCheck = regexp.MustCompile(`[A-Za-z0-9]+([\-\.]{1}[A-Za-z0-9]+)*`)
	return hostCheck.MatchString(host)
}

func main() {
	host := "http://localhost:8200/"
	//host := "http://10.218.9.169:8200/"
	u, e := url.Parse(host)
	if e != nil {
		fmt.Println("Error: ", e.Error())
	}
	fmt.Println(u.Host, u.Hostname(), u.Port())

	fmt.Println("validateVaultHost: ", validateVaultHost(u.Hostname()))
	addr := net.ParseIP(u.Hostname())

	if addr != nil {
		fmt.Println("this is valid IP ", u.Hostname())
	} else {
		fmt.Println("invalid iP: ", u.Hostname())
		fmt.Println("checking if its a hostname ", u.Hostname())
		hostName, err := net.LookupHost(u.Hostname())
		fmt.Println(hostName, err)
	}

}
