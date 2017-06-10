package wrapper

import (
	"bufio"
	"encoding/xml"
	"errors"
	"github.com/flamurey/tranzaq-connector/connector/logger"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	NotStarted int = iota
	Running
	Stopped
)

const (
	label_RequstResponce = "!@#$%"
	label_BrokerData     = "%$#@!"
)

var WrapperNotStarted = errors.New("Wrapper not started")
var WrapperIsRunningAlready = errors.New("Wrapper is running")
var WrapperClosed = errors.New("Wrapper closed")

//TranzaqWrapper represent a external process which work with broker tranzaq server
type TranzaqWrapper struct {
	Errors chan error

	cmd        exec.Cmd
	status     *int
	commands   chan tCommand
	answers    chan string
	started    *sync.WaitGroup
	brokerData chan string
}

func New(wrapperDir, winePath string) TranzaqWrapper {
	wrapperStderr := CreateWrapperOut()
	notStarted := NotStarted
	return TranzaqWrapper{
		cmd: exec.Cmd{
			Path:   winePath,
			Dir:    wrapperDir,
			Args:   []string{"wine", "wrapper.exe"},
			Stderr: wrapperStderr,
		},
		status:     &notStarted,
		commands:   make(chan tCommand),
		answers:    make(chan string),
		Errors:     make(chan error, 1000),
		started:    new(sync.WaitGroup),
		brokerData: make(chan string, 1000),
	}
}

func (wrapper TranzaqWrapper) Start() error {
	wrapper.started.Add(1)
	switch *wrapper.status {
	case Running:
		return WrapperIsRunningAlready
	case Stopped:
		return WrapperClosed
	}

	wrapperOut, err := wrapper.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	wrapperIn, err := wrapper.cmd.StdinPipe()
	if err != nil {
		return err
	}

	go wrapper.listenAnswers(wrapperOut)
	go wrapper.writeCommands(wrapperIn)
	go wrapper.processBrokerData()

	err = wrapper.cmd.Start()
	if err != nil {
		*wrapper.status = Stopped
		logger.WithError(err).Error("Cannot start wine app")
		return err
	}

	wrapper.started.Wait()
	*wrapper.status = Running
	logger.Info("Wrapper successful started")
	return nil
}

// logLevel : 1 -- minimal
//	2 -- standard
//	3 -- maximum
func (wrapper TranzaqWrapper) InitServer(logLevel int) error {
	_, err := wrapper.sendCommand(newInitServerCommand(logLevel))
	return err
}

func (wrapper TranzaqWrapper) ChangeLogLevel(logLevel int) error {
	_, err := wrapper.sendCommand(newChangeLogLevelCommand(logLevel))
	return err
}

func (wrapper TranzaqWrapper) DestroyServer() error {
	_, err := wrapper.sendCommand(newDestroyServerCommand())
	return err
}

func (wrapper TranzaqWrapper) SendCommand(command string) (string, error) {
	return wrapper.sendCommand(tCommand{
		command_Request,
		command,
		make(chan string),
	})
}

func (wrapper TranzaqWrapper) sendCommand(command tCommand) (answer string, err error) {
	wrapper.commands <- command
	answer = <-command.Answer
	switch command.commandType {
	case command_Request:
		response := new(tRequestResponse)
		parseErr := xml.Unmarshal([]byte(answer), response)
		if parseErr != nil {
			err = errors.New(answer)
		} else if !response.Success {
			err = errors.New(response.Message)
		}
	default:
		if answer != "" {
			err = errors.New(answer)
		}
	}
	if err != nil {
		logger.WithFields(logrus.Fields{
			"answer":  answer,
			"command": command.String(),
		}).Warn("Error perform command")
	}
	return
}

func (wrapper TranzaqWrapper) Close() error {
	logger.Info("Closing wrapper")
	switch *wrapper.status {
	case NotStarted:
		logger.WithError(WrapperNotStarted).Error("Illegal state")
		return WrapperNotStarted
	case Stopped:
		logger.WithError(WrapperClosed).Error("Illegal state")
		return WrapperClosed
	case Running:
		return wrapper.closeRunning()
	}
	return nil
}

func (wrapper TranzaqWrapper) closeRunning() error {
	_, err := wrapper.sendCommand(newCloseCommand())
	if err != nil {
		wrapper.Errors <- err
	}
	wrapper.cmd.Wait()
	return err
}

func (wrapper TranzaqWrapper) processBrokerData() {
	for data := range wrapper.brokerData {
		logger.Debug(data)
	}
}

func (wrapper TranzaqWrapper) listenAnswers(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			logger.Info("End read from wine wrapper")
			break
		} else if err != nil {
			logger.WithError(err).Error("Failed read input form wine wrapper")
			break
		}

		size := len(line)
		endIndex := size - 1
		//remove /r symbol if exists
		if line[size-2] == 13 {
			endIndex = size - 2
		}
		line = line[:endIndex]

		if strings.HasPrefix(line, label_RequstResponce) {
			flagIndex := len(label_RequstResponce)
			flag, _ := strconv.Atoi(line[flagIndex : flagIndex+1])
			if flag == 1 {
				wrapper.answers <- line[len(label_RequstResponce)+1:]
			} else {
				wrapper.answers <- ""
			}
		} else if strings.HasPrefix(line, label_BrokerData) {
			wrapper.brokerData <- line[len(label_BrokerData):]
		} else if line == "server started" {
			wrapper.started.Done()
		}
	}
}

func (wrapper TranzaqWrapper) writeCommands(writer io.Writer) {
	for command := range wrapper.commands {
		logger.WithField("command", command.String()).Info("Send new command to wrapper")
		_, err := writer.Write([]byte(command.String()))
		if err != nil {
			logger.WithError(err).WithField("command", command.String()).Error("Error on write command")
			wrapper.Errors <- err
			errAnswer, _ := xml.Marshal(tRequestResponse{
				Success: false,
				Message: err.Error(),
			})
			command.Answer <- string(errAnswer)
		} else {
			answer := <-wrapper.answers
			command.Answer <- answer
		}
		close(command.Answer)
	}
}
