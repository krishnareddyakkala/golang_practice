package main

// this works in golang 1.15 and above only

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/ocsp"
)

func main() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	tlsConfig.VerifyConnection = func(cs tls.ConnectionState) error {
		fmt.Println(len(cs.PeerCertificates))
		r, e := ocsp.ParseResponse(cs.OCSPResponse, cs.PeerCertificates[1])
		if e != nil {
			return nil
		}
		fmt.Println(r.Status)
		if r.Status == 1 {
			return errors.New("certificate revoked")
		}
		return nil
	}
	req, _ := http.NewRequest("GET", "http://tstdmzihealth-api.olympus.f5net.com/", nil)
	cl := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}
	resp, err := cl.Do(req)
	if err != nil {
		fmt.Println("error occurred 4", err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("###########", string(body))
}
