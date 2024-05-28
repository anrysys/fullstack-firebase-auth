// global/config.go
package global

import (
	"backend/configs"
	"log"
)

var Conf *configs.Config

func LoadConfig() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	Conf = &config
}
