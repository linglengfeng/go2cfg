package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	Project *viper.Viper
)

func init() {
	projectConfigJson()
	fmt.Println("project config all init successed.")
}

func projectConfigJson() {
	Project = viper.New()
	Project.AddConfigPath("./")
	Project.SetConfigName("config.json")
	Project.SetConfigType("json")
	if err := Project.ReadInConfig(); err != nil {
		fmt.Println("projectConfigJson err:", err)
	}
}
