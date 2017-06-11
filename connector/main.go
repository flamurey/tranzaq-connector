package main

import (
	//"flag"
	//pb "github.com/flamurey/prototrade"
	//"golang.org/x/net/context"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/grpclog"
	//"github.com/flamurey/tranzaq-connector/connector/wrapper"
	//"log"
	//"os"

	"io"
	"strconv"
	"strings"

	//"context"
	"flag"
	"github.com/chzyer/readline"
	log "github.com/flamurey/tranzaq-connector/connector/logger"
	"github.com/flamurey/tranzaq-connector/connector/wrapper"
)

const (
	ConnectCommand    = "connect"
	RequestCommand    = "command"
	LogChangeCommand  = "cange_log_level"
	DisconnectCommand = "disconnect"
	ExistCommand      = "exit"
)

var (
	wrapperDir = flag.String("wrapper_path", "/home/lesha/ws/ats/tranzaq-docker/bin", "path to exe of wrapper program")
	winePath   = flag.String("wine_path", "/usr/bin/wine", "The path to wine for runnig tranzaq wrapper")
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem(ConnectCommand),
	readline.PcItem(RequestCommand),
	readline.PcItem(LogChangeCommand),
	readline.PcItem(DisconnectCommand),
	readline.PcItem(ExistCommand),
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

func main() {
	flag.Parse()
	log.Init()
	//start tranzaq wrapper
	tranzaq := wrapper.New(*wrapperDir, *winePath)
	err := tranzaq.Start()
	if err != nil {
		log.WithError(err).Fatalf("Wrapper program couldn't start")
	}
	defer tranzaq.Close()

	l, err := readline.NewEx(&readline.Config{
		Prompt:            "\033[31m»\033[0m ",
		HistoryFile:       "/tmp/tranzaq-connector.tmp",
		AutoComplete:      completer,
		HistorySearchFold: true,
	})
	if err != nil {
		log.WithError(err).Fatalf("Fails to start readline")
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {

		case strings.HasPrefix(line, ConnectCommand):
			commandArg := strings.TrimSpace(line[len(ConnectCommand):])
			logLevel, err := strconv.Atoi(commandArg)
			if err != nil {
				log.WithField("argument", commandArg).Warn("Fail to parse logLevel argument to int")
			} else {
				tranzaq.InitServer(logLevel)
			}

		case strings.HasPrefix(line, RequestCommand):
			command := strings.TrimSpace(line[len(RequestCommand):])
			tranzaq.SendCommand(command)

		case strings.HasPrefix(line, LogChangeCommand):
			commandArg := strings.TrimSpace(line[len(LogChangeCommand):])
			logLevel, err := strconv.Atoi(commandArg)
			if err != nil {
				log.WithField("argument", commandArg).Warn("Fail to parse logLevel argument to int")
			} else {
				tranzaq.ChangeLogLevel(logLevel)
			}

		case strings.HasPrefix(line, DisconnectCommand):
			tranzaq.DestroyServer()

		case line == ExistCommand:
			goto exit
		case line == "help":
			usage(l.Stderr())
		default:
			usage(l.Stderr())
		}
	}
exit:
}
