package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	r := regexp.MustCompile(`((https:\/\/)+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}))`) // https://1.1.1.1
	//r := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`) // 1.0.0.1
	err := filepath.Walk("query",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			dat, err1 := ioutil.ReadFile(path)
			if err1 != nil {
				fmt.Println("$$$$$$$$$$$$ cant open file: ", path)
				return nil
			}
			res := r.ReplaceAll(dat, []byte("0.0.0.0"))
			er := ioutil.WriteFile(path, res, 0766)
			if er != nil {
				fmt.Printf("###### file: %s, error: %s", path, er.Error())
			}
			fmt.Println("@@@@@@@ content work done")
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}

}
