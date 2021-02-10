package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Records struct {
	Name string
	Id   int
}
type FinalDetails struct {
	Records interface{}
	Name    string
	Rec     []Records
}

//
//type FinalReverse struct {
//}

type NestRecords struct {
	Records
	NestVal int
}

func main() {

	records := []Records{
		{
			Name: "Name 1",
			Id:   1,
		},
		{
			Name: "Name 2",
			Id:   2,
		},
	}
	recordBytes, _ := json.Marshal(records)
	fmt.Println(string(recordBytes))
	finalRecords := FinalDetails{Name: "some Name"}
	finalRecords.Rec = records
	finalRecords.Records = "####"
	finalBytes, _ := json.Marshal(finalRecords)
	finalStr := string(finalBytes)
	recordStr := string(recordBytes)
	recordStr = strings.Trim(recordStr, "\"")
	finalStr = strings.Replace(finalStr, "\"####\"", strings.Trim(recordStr, "\""), 1)
	fmt.Println(finalStr)

	var finalReverse FinalDetails

	_ = json.Unmarshal([]byte(finalStr), &finalReverse)
	finalRevBytes, _ := json.MarshalIndent(finalReverse, "", "\t")
	fmt.Println(string(finalRevBytes))

}
