package console

import (
	"github.com/chzyer/readline"
	log "github.com/flamurey/tranzaq-connector/connector/logger"
	"github.com/flamurey/tranzaq-connector/connector/wrapper"
	"io"
	"strconv"
	"strings"
)

const (
	ConnectCommand    = "connect"
	RequestCommand    = "command"
	LogChangeCommand  = "cange_log_level"
	DisconnectCommand = "disconnect"
	ExistCommand      = "exit"
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem(ConnectCommand),
	readline.PcItem(RequestCommand),
	readline.PcItem(LogChangeCommand),
	readline.PcItem(DisconnectCommand),
	readline.PcItem(ExistCommand),
)

func StartReadInput(tranzaq wrapper.TranzaqWrapper) {
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
