package controller

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

//Configuration struct
type Configuration struct {
	Port             string // port no
	ConnectionString string // connection string
	Database         string // database name
	Collection       string // collection
}

/*ReadConfig Reading the configs from  db.properties
 */
func ReadConfig() Configuration {
	var configfile = "config.properties"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Configuration
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func HandleRequests() {

	config := ReadConfig()
	var port = ":" + config.Port

	r := gin.Default()
	r.GET("/", Homepage)
	r.POST("/courses", Addcoures)
	r.GET("/courses", FindAllCources)
	r.DELETE("/courses/:id", DeleteCourse)
	r.PUT("/courses/:id", UpdateCourse)
	r.POST("/students", CreateStudent)
	r.GET("/students", Getstudent)
	// r.GET("/students", Aggregate)

	r.Run()
	fmt.Printf("application listening port%s\n", port)

}
