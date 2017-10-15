package commands

import (
	"encoding/xml"
)

type Command struct {
	XMLName xml.Name `xml:"command"`
	Id      string   `xml:"id,attr"`
}

func (c *Command) Name() string {
	return c.Id
}

func NewDisconnect() Command {
	return Command{Id: "disconnect"}
}

func NewServerStatus() Command {
	return Command{Id: "server_status"}
}

func NewConnect() Command {
	return Command{Id: "connect"}
}

type ConnectCommand struct {
	Command
	Login          string `xml:"login"`
	Password       string `xml:"password"`
	Host           string `xml:"host"`
	Port           int    `xml:"port"`
	Autopos        bool   `xml:"autopos"`          // true by default into wrapper
	RequestDelay   int    `xml:"rqdelay"`          // 100mc min for default version of Tranzaq and 10mc for HFT
	SessionTimeout int    `xml:"session_timeout"`  // 120 sec by default
	RequestTimeout int    `xml:"srequest_timeout"` // 20 sec by default
	Milliseconds   bool   `xml:"milliseconds"`
	UtcTime        bool   `xml:"utc_time"`
}
