package gophast

import "log"

type ConfigT struct {
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

func (config *ConfigT) Dump() {
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
		log.Printf("API Host:			 %s", config.APIHost)
	}
	log.Printf("GoPhast Port:		 %d", config.Port)
}
