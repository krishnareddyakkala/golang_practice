package main

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/crypto/ocsp"
)

func loadCert(path string) (*x509.Certificate, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	bytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	_ = block
	// return x509.ParseCertificate(bytes)
	return x509.ParseCertificate(bytes)
}

func main() {
	cert, err := loadCert(os.Args[1])
	if err != nil {
		panic(err)
	}
	ca, err := loadCert(os.Args[2])
	if err != nil {
		panic(err)
	}

	revoked := isCertificateRevokedByOCSP(cert, ca, "http://ocsp.disa.mil")
	fmt.Printf("Revoked: %t\n", revoked)
}

func isCertificateRevokedByOCSP(clientCert, issuerCert *x509.Certificate, ocspServer string) bool {
	opts := &ocsp.RequestOptions{Hash: crypto.SHA1}
	buffer, err := ocsp.CreateRequest(clientCert, issuerCert, opts)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return false
	}
	httpRequest, err := http.NewRequest(http.MethodPost, ocspServer, bytes.NewBuffer(buffer))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return false
	}
	ocspUrl, err := url.Parse(ocspServer)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return false
	}
	httpRequest.Header.Add("Content-Type", "application/ocsp-request")
	httpRequest.Header.Add("Accept", "application/ocsp-response")
	httpRequest.Header.Add("host", ocspUrl.Host)
	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return false
	}
	defer httpResponse.Body.Close()
	output, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return false
	}

	ocspResponse, err := ocsp.ParseResponse(output, issuerCert)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return false
	}
	if ocspResponse.Status == ocsp.Revoked {
		log.Printf("certificate has been revoked by OCSP server %s, refusing connection", ocspServer)
		return true
	} else {
		return false
	}
}
