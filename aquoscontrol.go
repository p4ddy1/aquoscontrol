package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/tarm/serial"
)

type Config struct {
	Logging  bool
	Server   serverConfig
	Serial   serialConfig
	Commands map[string]cmd
}

type serverConfig struct {
	Address string
	Path    string
}

type serialConfig struct {
	Port string
	Baud int
}

type cmd struct {
	Command      string
	HasParameter bool
}

type ControlHandler struct {
	config *serial.Config
}

var conf Config
var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.toml", "-config=/path/to/config.toml")
}

func main() {
	flag.Parse()

	_, err := toml.DecodeFile(configPath, &conf)
	if err != nil {
		writeToLog(err, true)
	}

	serialConfig := &serial.Config{Name: conf.Serial.Port, Baud: conf.Serial.Baud}

	ch := &ControlHandler{config: serialConfig}
	http.Handle(conf.Server.Path, ch)
	err = http.ListenAndServe(conf.Server.Address, nil)
	if err != nil {
		writeToLog(err, true)
	}
}

func generateCommand(cmd, parameter string) string {
	command := cmd + parameter
	return command + strings.Repeat(" ", 8-len(command))
}

func writeToLog(text interface{}, fatal bool) {
	if !fatal {
		if conf.Logging {
			log.Println(text)
		}
	} else {
		log.Fatal(text)
	}
}

func (ch *ControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cmdReq := r.FormValue("cmd")
		value := r.FormValue("value")
		if command, res := conf.Commands[cmdReq]; res {
			if command.HasParameter {
				ch.WritePort(generateCommand(command.Command, value))
			} else {
				ch.WritePort(generateCommand(command.Command, ""))
			}
		} else {
			writeToLog("Command not found: "+cmdReq, false)
		}
	}
}

func (ch *ControlHandler) WritePort(cmd string) {
	s, err := serial.OpenPort(ch.config)
	if err != nil {
		writeToLog(err, false)
		return
	}
	_, err = s.Write([]byte(cmd + "\r\n"))
	if err != nil {
		writeToLog(err, false)
		return
	}
	s.Close()
	writeToLog("Executed command: "+cmd, false)
}
