package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var config configT

// var publicDir = ""
// var subsite = ""
// var assetRoute = ""
// var apiRoute = ""
// var apihost = ""

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func handlePages(c *gin.Context) {
	part1 := c.Param("part1")
	part2 := c.Param("part2")

	log.Printf("Got request %s, part1 = %s, part2 = %s", c.Request.URL.Path, part1, part2)

	if config.IsAPI && config.APIHost != "" && part1 == config.APIRoute {
		log.Printf("Forwarding through handlePages router")
		forwardAPIRequest(c)
	}

	log.Printf("config.AssetRoute: %s", config.AssetRoute)
	if part1 != config.AssetRoute {
		log.Printf("Sending landing page")
		getLandingPage(c)
	}

	if part1 == config.AssetRoute && part2 != "" {
		log.Printf("Sending asset")
		getPublicAsset(c)
	}
}

func forwardAPIRequest(c *gin.Context) {
	log.Printf("Forwarding api request.")

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = config.APIHost

		newURL := req.URL.String()
		log.Printf("Replacing url header...")
		newURL = strings.Replace(newURL, config.APIRoute+"/", "", -1)
		log.Printf("Got new url: %s\n", newURL)

		url, err := url.Parse(newURL)
		log.Printf("URL: %s, ERR: %s\n", url.String(), err)
		if url.String() != "" && err == nil {
			log.Printf("Forwarding request from '%s' to '%s'\n",
				req.URL, newURL)
			req.URL = url
		} else {
			log.Printf("Could not forward request from '%s' to '%s'\n",
				req.URL, newURL)
		}
	}

	log.Printf("Successfully unwrapped URL")
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func getLandingPage(c *gin.Context) {
	landingName := c.Param("part1")
	subPageName := c.Param("part2")
	if landingName != "" && subPageName != "" {
		// log.Printf("Trying to get subpage, landing = %s, subpage = %s", landingName, subPageName)
		landingName = landingName + "/" + subPageName
	}

	if landingName == "" {
		landingName = "index"
	}

	if config.IsSubsite && config.Subsite != "" {
		landingName = config.Subsite + "/" + landingName
	}
	// log.Printf("Landing page requested, %s", landingName)

	// log.Printf("... Looking at path %s", "./public/html/" + landingName + ".html")
	c.File(config.PublicDir + "/html/" + landingName + ".html")
}

func getPublicAsset(c *gin.Context) {
	imageTypes := []string{"jpg", "png", "tiff", "jpeg", "gif", "svg"}

	fileName := c.Param("part2")
	fileParts := strings.Split(fileName, ".")
	fileType := fileParts[len(fileParts)-1]
	// log.Printf("Getting file %s", fileName)

	if contains(imageTypes, fileType) {
		// log.Printf("Requested fileType %s is an image", fileType)
		fileType = "img"
	}

	c.File(config.PublicDir + "/" + fileType + "/" + fileName)
}

type yamlT struct {
	Name    string `yaml:"name"`
	Port    int32  `yaml:"port"`
	Release bool   `yaml:"release"`

	PublicDir  string `yaml:"public"`
	AssetRoute string `yaml:"asset_route"`

	Subsite string `yaml:"subsite"`

	APIRoute string `yaml:"api_route"`
	APIHost  string `yaml:"api_host"`
}

type configT struct {
	IsNamed bool
	Name    string
	Port    int32
	Release bool

	PublicDir  string
	AssetRoute string

	Subsite   string
	IsSubsite bool

	IsAPI    bool
	APIRoute string
	APIHost  string
}

func (c *yamlT) getConfigData() *yamlT {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	var yamlData yamlT
	yamlData.getConfigData()

	// BEGIN: Config based on YAML data
	if yamlData.Name != "" {
		config.IsNamed = true
		config.Name = yamlData.Name
	}

	if yamlData.Port == 0 {
		config.Port = 4000
	} else {
		config.Port = yamlData.Port
	}

	config.Release = yamlData.Release
	if yamlData.Release {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if yamlData.PublicDir != "" {
		config.PublicDir = yamlData.PublicDir
	} else {
		config.PublicDir = "./public"
	}

	if yamlData.AssetRoute != "" {
		config.AssetRoute = yamlData.AssetRoute
	} else {
		config.AssetRoute = "a"
	}

	if yamlData.Subsite != "" {
		config.Subsite = "_" + yamlData.Subsite
		config.IsSubsite = true
	}

	if yamlData.APIRoute != "" && yamlData.APIHost != "" {
		config.APIRoute = yamlData.APIRoute
		config.APIHost = yamlData.APIHost
		config.IsAPI = true
	}
	// END

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // Default, send 500 on panic

	r.GET("/", handlePages)
	r.GET("/:part1", handlePages)
	r.GET("/:part1/:part2", handlePages)
	r.GET("/:part1/:part2/*part3", handlePages)

	if config.IsAPI {
		r.POST("/"+config.APIRoute+"/*part1", forwardAPIRequest)
	}

	// Listen and serve on 0.0.0.0:8080
	// BEGIN: Dump config data
	if config.IsNamed {
		log.Printf("Server: 			 %s", config.Name)
	}
	if config.IsSubsite {
		log.Printf("Subsite: 			 %s", config.Subsite)
	}
	if config.Release {
		log.Printf("Environment:		 release")
	}
	log.Printf("Public Dir:			 %s", config.PublicDir)
	if config.IsAPI {
		log.Printf("API Path:			 /%s", config.APIRoute)
	}
	log.Printf("Port:			 %d", config.Port)
	// END
	r.Run(":" + fmt.Sprintf("%d", config.Port))
}
