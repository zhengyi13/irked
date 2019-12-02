package main

import (
	"log"
	"net"
	"strconv"

	"flag"

	"io/ioutil"

	"gopkg.in/irc.v3"
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

	// TODO: move this out to a parseConfig() function. Try to make it
	// package "irked", move this file also to package "irked".
	configfile := flag.String("config", "", "path to your config file")

	flag.Parse()
	confdata, err := ioutil.ReadFile(*configfile)
	if err != nil {
		log.Fatalln(err)
	}
	irc_server_config := ServerList{}
	if err = yaml.Unmarshal(confdata, &irc_server_config); err != nil {
		log.Print(err)
		log.Fatalln("Died trying to unmarshall config from file")
	}

	// TODO: loop through servers, prefer a config where Tls is true, then connect to that
	// using tls.Dial

	hostport := irc_server_config.Servers[0].Host + ":" + strconv.Itoa(irc_server_config.Servers[0].Port)

	conn, err := net.Dial("tcp", hostport)
	if err != nil {
		log.Fatalln(err)
	}

	config := irc.ClientConfig{
		Nick: irc_server_config.Servers[0].Nick,
		Pass: irc_server_config.Servers[0].Password,
		User: irc_server_config.Servers[0].Username,
		Name: irc_server_config.Servers[0].Name,
		Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
			if m.Command == "001" {
				// 001 is a welcome event, so we join channels there
				c.Write("JOIN #bot-test-chan")
			} else if m.Command == "PRIVMSG" && c.FromChannel(m) {
				// Create a handler on all messages.
				c.WriteMessage(&irc.Message{
					Command: "PRIVMSG",
					Params: []string{
						m.Params[0],
						m.Trailing(),
					},
				})
			}
		}),
	}

	// Create the client
	client := irc.NewClient(conn, config)
	err = client.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
