package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
)

func main() {
	fmt.Println(getTLSConfig())

}

func getTLSConfig() (*tls.Config, error) {

	rootCAs, _ := x509.SystemCertPool()
	for _, c := range rootCAs.Subjects() {
		if strings.Contains(string(c), "Entrust, Inc.") {
			fmt.Println("subject: ", string(c))
		}
	}
	//if rootCAs == nil {
	//	rootCAs = x509.NewCertPool()
	//}
	//cacert, err := ioutil.ReadFile("/certs/ca.crt")
	//if err != nil {
	//	return nil, errors.New("could not read CA cert file")
	//}
	//rootCAs.AppendCertsFromPEM(cacert)
	//tlsConfig := &tls.Config{
	//	RootCAs:   rootCAs,
	//	ClientCAs: rootCAs,
	//}
	//cert, err := tls.LoadX509KeyPair("/certs/client.crt", "/certs/client.key")
	//if err != nil {
	//	return nil, errors.New("could not create TLS X509 keypair")
	//}
	//
	//tlsConfig.Certificates = []tls.Certificate{cert}
	//
	//return tlsConfig, nil
	return &tls.Config{}, nil
}
