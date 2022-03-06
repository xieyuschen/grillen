package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var (
	port *string
)

func SettingUpEnvironment() {
	c := conf.ReadSettingsFromFile("Config.json")
	initArgs(c.Version)
	dbclt.InitDb(c.DbSettings)
}
func main() {
	SettingUpEnvironment()
	r := gin.Default()
	router.UseHoleRouter(r)
	des := ":" + *port
	_ = r.Run(des)
}

func initArgs(version string) {
	port = flag.String("port", "8082", "Listen port")
	flag.Parse()

	//Check whether json file version is match to server version,avoid using
	//Develop Server json file on deployed server
	if *port != version {
		panic("Input json doesn't match server!! Pay attention to its version")
	}
}
