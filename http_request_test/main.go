package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	url := "http://10.218.9.169:8200/ui/vault/secrets/shared-service-secrets/list"
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		fmt.Println(err)
		return
	}
	req = req.WithContext(ctx)

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{}}}
	res, err := httpClient.Do(req)
	//fmt.Println(res.Status, res.StatusCode)

	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

}
