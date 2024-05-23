// global/config.go
package global

import (
	"backend/configs"
	"log"
)

var Conf *configs.Config

func LoadConfig(path string) {
	config, err := configs.LoadConfig(path)
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	Conf = &config
}
