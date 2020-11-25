package gophast

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var Config ConfigT

func HandlePages(c *gin.Context) {
	part1 := c.Param("part1")
	part2 := c.Param("part2")

	log.Printf("Got request %s, part1 = %s, part2 = %s", c.Request.URL.Path, part1, part2)

	if Config.IsAPI && Config.APIHost != "" && part1 == Config.APIRoute {
		log.Printf("Forwarding through handlePages router")
		ForwardAPIRequest(c)
	} else if part1 != Config.AssetRoute {
		log.Printf("Sending landing page")
		GetLandingPage(c)
	} else if part1 == Config.AssetRoute && part2 != "" {
		log.Printf("Sending asset")
		GetPublicAsset(c)
	}
}

func ForwardAPIRequest(c *gin.Context) {
	log.Printf("Forwarding api request.")

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = Config.APIHost

		newURL := req.URL.String()
		log.Printf("Replacing url header...")
		newURL = strings.Replace(newURL, Config.APIRoute+"/", "", -1)
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

func GetLandingPage(c *gin.Context) {
	landingName := c.Param("part1")
	subPageName := c.Param("part2")
	if landingName != "" && subPageName != "" {
		// log.Printf("Trying to get subpage, landing = %s, subpage = %s", landingName, subPageName)
		landingName = landingName + "/" + subPageName
	}

	if landingName == "" {
		landingName = "index"
	}

	if Config.IsSubsite && Config.Subsite != "" {
		landingName = Config.Subsite + "/" + landingName
	}
	// log.Printf("Landing page requested, %s", landingName)

	// log.Printf("... Looking at path %s", "./public/html/" + landingName + ".html")
	c.File(Config.PublicDir + "/html/" + landingName + ".html")
}

func GetPublicAsset(c *gin.Context) {
	imageTypes := []string{"jpg", "png", "tiff", "jpeg", "gif", "svg"}

	fileName := c.Param("part2")
	fileParts := strings.Split(fileName, ".")
	fileType := fileParts[len(fileParts)-1]
	// log.Printf("Getting file %s", fileName)

	if Contains(imageTypes, fileType) {
		// log.Printf("Requested fileType %s is an image", fileType)
		fileType = "img"
	}

	c.File(Config.PublicDir + "/" + fileType + "/" + fileName)
}
