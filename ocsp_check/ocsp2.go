package main

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/crypto/ocsp"
)

func main() {

	// Issue an HTTP GET request to a server. `http.Get` is a
	// convenient shortcut around creating an `http.Client`
	// object and calling its `Get` method; it uses the
	// `http.DefaultClient` object which has useful default
	// settings.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	tlsConfig.VerifyPeerCertificate = func(certificates [][]byte, _ [][]*x509.Certificate) error {
		certs := make([]*x509.Certificate, len(certificates))
		for i, asn1Data := range certificates {
			cert, err := x509.ParseCertificate(asn1Data)
			if err != nil {
				return errors.New("tls: failed to parse certificate from server: " + err.Error())
			}
			certs[i] = cert
		}
		opts := x509.VerifyOptions{
			Roots:         tlsConfig.RootCAs, // On the server side, use config.ClientCAs.
			DNSName:       tlsConfig.ServerName,
			Intermediates: x509.NewCertPool(),
			// On the server side, set KeyUsages to ExtKeyUsageClientAuth. The
			// default value is appropriate for clients side verification.
			// KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		}
		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}
		fmt.Println(opts.CurrentTime)
		//ocspResponse, err := ocsp.CreateResponse(certs[0], certs[1],nil,certs[0].PublicKey)
		//if err != nil {
		//    fmt.Println("%%%%%%%%%%%", err)
		//    return err
		//}
		//fmt.Println("^^^^^^^^^^^", ocspResponse)
		if isCertificateRevokedByOCSP2(certs[0].Issuer.CommonName, certs[0], certs[1], certs[0].OCSPServer[0]) {
			return errors.New("certificate revoked")
		}
		//_, err := certs[0].Verify(opts)
		return nil
	}
	req, _ := http.NewRequest("GET", "https://gorest.co.in/public-api/users", nil)
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

func isCertificateRevokedByOCSP2(commonName string, clientCert, issuerCert *x509.Certificate, ocspServer string) bool {
	opts := &ocsp.RequestOptions{Hash: crypto.SHA1}
	buffer, err := ocsp.CreateRequest(clientCert, issuerCert, opts)
	if err != nil {
		return false
	}
	httpRequest, err := http.NewRequest(http.MethodPost, ocspServer, bytes.NewBuffer(buffer))
	if err != nil {
		return false
	}
	ocspUrl, err := url.Parse(ocspServer)
	if err != nil {
		return false
	}
	httpRequest.Header.Add("Content-Type", "application/ocsp-request")
	httpRequest.Header.Add("Accept", "application/ocsp-response")
	httpRequest.Header.Add("host", ocspUrl.Host)
	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return false
	}
	defer httpResponse.Body.Close()
	output, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		fmt.Println("error while reading body", err.Error())
		return false
	}
	ocspResponse, err := ocsp.ParseResponse(output, issuerCert)
	if err != nil {
		fmt.Println("error while ocsp response", err.Error())
		return false
	}
	if ocspResponse.Status == ocsp.Revoked {
		fmt.Printf("certificate '%s' has been revoked by OCSP server %s, refusing connection\n", commonName, ocspServer)
		return true
	} else {
		return false
	}
}
