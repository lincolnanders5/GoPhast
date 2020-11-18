package gophast

import "github.com/gin-gonic/gin"

func InitConfigVars(config *ConfigT, yamlData *YamlT) {
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
}