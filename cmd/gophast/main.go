package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lincolnanders5/GoPhast/internal/gophast"
	"log"
)

var config gophast.ConfigT

func main() {
	var yamlData gophast.YamlT
	yamlData.GetConfigData()

	// BEGIN: Config based on YAML data
	gophast.InitConfigVars(&config, &yamlData)
	// END

	gophast.Config = config

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // Default, send 500 on panic

	r.GET("/", gophast.HandlePages)
	r.GET("/:part1", gophast.HandlePages)
	r.GET("/:part1/*part2", gophast.HandlePages)

	if config.IsAPI {
		r.POST("/"+config.APIRoute+"/*part1", gophast.ForwardAPIRequest)
	}

	// Listen and serve on 0.0.0.0:8080
	// BEGIN: Dump config data
	config.Dump()
	// END
	err := r.Run(":" + fmt.Sprintf("%d", config.Port))
	log.Fatal(err)
}