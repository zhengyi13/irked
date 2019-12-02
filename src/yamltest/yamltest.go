package main

// An attempt to re-learn YAML parsing, pursuant to designing a YAML-based
// config for irked.go

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ServerList struct {
	Servers []ServerConfig `yaml:"servers"`
}

type ServerConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Tls      bool   `yaml:"tls"`
	Nick     string `yaml:"nick"`
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
	Fullname string `yaml:"fullname"`
}

func main() {
	data, err := ioutil.ReadFile("./files/myconfig.yml")
	if err != nil {
		log.Fatalln(err)
	}

	irc_config := ServerList{}
	if err = yaml.Unmarshal(data, &irc_config); err != nil {
		log.Print(err)
		log.Fatalln("Died trying to unmarshall my config")
	}

	for i, server := range irc_config.Servers {
		log.Printf("Server #%v: \n", i)
		log.Printf("%v:%v, TLS: %v", server.Host, server.Port, server.Tls)
	}
}
