package main

import (
	"fmt"
	"runtime"

	"github.com/rdscfh/spider/Node"
	"github.com/spf13/viper"
)

func main() {
	var MULTICORE int = runtime.NumCPU() //number of core
	runtime.GOMAXPROCS(MULTICORE)        //running in multicore

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error when reading config file: %s\n", err))
	}

	url := viper.GetString("url")

	spider.Run(url)
}
