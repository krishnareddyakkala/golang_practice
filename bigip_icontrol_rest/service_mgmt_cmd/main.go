// parts of this code copied from `deos-agent` repo
// for better understanding of Managing services in BigIP with iControl REST
// refer to this https://support.f5.com/csp/article/K51226856 link

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BigIPDetails struct {
	Version      string `db:"version"`
	IPAddress    string `db:"ip_address"`
	CredUser     string `db:"cred_user"`
	CredPassword string `db:"cred_passwd"`
}

type DataSet struct {
	Name    string `json:"name,omitempty"`
	Url     string `json:"url,omitempty"`
	Method  string `json:"method,omitempty"`
	ReqBody string `json:"req_body,omitempty"`
}

type BasicAuth struct {
	Username string
	Password string
}

func runServiceCmd(producer BigIPDetails, dataset DataSet) (data interface{}, err error) {
	requestContext := prepHttpRequestContext(producer, dataset)

	body, responseCode, err := executeHTTPRequest(requestContext)
	if err != nil {
		fmt.Printf("error getting data from big-ip %s error is %v  ", dataset.Name, err)
		return nil, err
	}
	if responseCode.StatusCode != http.StatusOK {
		fmt.Printf("non HTTP 200 response from  big-ip producer %s  - dataset name - %s http response %d ", producer.IPAddress, dataset, responseCode.StatusCode)
		return data, nil
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("data is an invalid JSON", err)
		return data, err
	}
	var prettyJSON bytes.Buffer
	//prettying JSON to for better debugging
	fmt.Printf("http response status for %s  is %s", dataset.Name, responseCode.Status)
	_ = json.Indent(&prettyJSON, body, "", "\t")
	fmt.Printf("%s raw data from BIG-IP %s - size is %s", dataset.Name, producer.IPAddress, prettyJSON.String())
	return data, nil
}

func getServiceMgmtDataSet(serviceName, command string) DataSet {
	dataset := DataSet{}
	dataset.Method = http.MethodPost
	dataset.Name = fmt.Sprintf("%s-%s", command, serviceName)
	dataset.Url = "/mgmt/tm/sys/service"
	dataset.ReqBody = fmt.Sprintf(`{"command":"%s","name":"%s"}`, command, serviceName)
	return dataset
}

// refer to this https://support.f5.com/csp/article/K51226856 for more details
func main() {

	details := BigIPDetails{}
	details.IPAddress = "https://**.**.**.**"
	details.CredUser = "admin"
	details.CredPassword = "***************"

	_, _ = runServiceCmd(details, getServiceMgmtDataSet("httpd", "restart"))
}
