package wrapper

import (
	"encoding/xml"
	"strconv"
)

type commandType uint8

const (
	command_InitServer    commandType = 4
	command_Request       commandType = 1
	command_ChangeLog     commandType = 2
	command_DestroyServer commandType = 3
	command_Close         commandType = 9
)

type tCommand struct {
	commandType
	Message string
	Answer  chan string
}

type tRequestResponse struct {
	XMLName xml.Name `xml:"result"`
	Success bool     `xml:"success,attr"`
	Message string   `xml:"message,omitempty"`
}

func (c tCommand) String() string {
	return strconv.Itoa(int(c.commandType)) + c.Message + "\n"
}

func newCloseCommand() tCommand {
	return tCommand{
		command_Close,
		"",
		make(chan string),
	}
}

func newInitServerCommand(logLevel int) tCommand {
	return tCommand{
		command_InitServer,
		strconv.Itoa(logLevel),
		make(chan string),
	}
}

func newChangeLogLevelCommand(logLevel int) tCommand {
	return tCommand{
		command_ChangeLog,
		strconv.Itoa(logLevel),
		make(chan string),
	}
}

func newDestroyServerCommand() tCommand {
	return tCommand{
		command_DestroyServer,
		"",
		make(chan string),
	}
}
