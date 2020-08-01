package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Port             int `env:"APP_PORT"`
	ConnectionString string
}

func main() {

	// mock env variable
	os.Setenv("ConnectionString", "test")
	//os.Setenv("APP_PORT", "8085")

	configuration := Configuration{}
	err := gonfig.GetConf("config/default.json", &configuration)
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}

	fmt.Println(configuration.Port)
	fmt.Println(configuration.ConnectionString)

	envName, iseEnvNameExists := os.LookupEnv("ENV")
	if iseEnvNameExists {
		fmt.Printf("$$$$$$$$$$$$$$ Reading '%s' environment config file $$$$$$$$$$$$$$\n", envName)
		err = gonfig.GetConf(getFileName(), &configuration)
		if err != nil {
			fmt.Println(err)
			os.Exit(500)
		}
		fmt.Println(configuration.Port)
		fmt.Println(configuration.ConnectionString)
	}

}

func getFileName() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		return "default"
	}
	filename := []string{"config/", env, ".json"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
