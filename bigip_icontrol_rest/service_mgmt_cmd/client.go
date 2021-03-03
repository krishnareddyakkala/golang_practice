package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type HTTPRequestContext struct {
	URL       string
	Method    string
	Header    http.Header
	Body      io.Reader
	Transport *http.Transport
	Auth      BasicAuth
}

func executeHTTPRequest(requestContext HTTPRequestContext) ([]byte, *http.Response, error) {
	fmt.Println("executing http requestContext ", requestContext)

	request, err := prepHttpRequest(requestContext)

	if err != nil {
		fmt.Println("error occurred preparing http.Request - ", err)
		return nil, nil, err
	}

	httpClient := &http.Client{Transport: requestContext.Transport}

	res, err := httpClient.Do(request)
	if err != nil {
		fmt.Println("error executing http.Request ", err)
		return nil, nil, err
	}

	defer Close(res.Body)

	body, _ := ioutil.ReadAll(res.Body)
	return body, res, nil
}

func Close(readCloser io.ReadCloser) {
	err := readCloser.Close()
	if err != nil {
		fmt.Println("Error in closing ReadCloser", err)
	}
}

func prepHttpRequest(requestContext HTTPRequestContext) (*http.Request, error) {

	request, err := http.NewRequest(requestContext.Method, requestContext.URL, requestContext.Body)

	if err != nil {
		return nil, err
	}

	request.Close = true
	if requestContext.Auth.Username != "" {
		request.SetBasicAuth(requestContext.Auth.Username, requestContext.Auth.Password)
	}

	for key, value := range requestContext.Header {
		request.Header.Add(key, value[0])
	}
	return request, nil
}

func prepHttpRequestContext(producer BigIPDetails, dataset DataSet) HTTPRequestContext {
	url := producer.IPAddress + dataset.Url
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	return HTTPRequestContext{
		URL:       url,
		Method:    dataset.Method,
		Header:    http.Header{},
		Body:      strings.NewReader(dataset.ReqBody),
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
		Auth:      BasicAuth{Username: producer.CredUser, Password: producer.CredPassword},
	}
}
