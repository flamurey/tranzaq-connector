package logger

import (
	"flag"
	"fmt"
	"github.com/flamurey/tranzaq-connector/connector/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

var (
	withLogstash    = flag.Bool("with_ls", true, "Send logs to logstash")
	logstashPort    = flag.Int("ls_port", 10001, "Logstash port")
	logstashAddress = flag.String("ls_ip", "localhost", "Logstash ip address")
	logLevel        = flag.String("log_level", "info", "Logs level filter, one of: panic, fatal, warn, info, debug")
)

var logEntry *logrus.Entry

var inited = false

func Init() {
	if inited {
		return
	}
	inited = true
	//flag.Parse()
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	if *withLogstash {
		serverAddr := fmt.Sprint(*logstashAddress, ":", *logstashPort)
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			Fatal("Cannot connect to logstah. Logging into log file.")
		}
		log.SetOutput(conn)
	} else {
		logFile, err := os.OpenFile("tranzaq.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			Fatalf("error opening file: %v", err)
		}
		log.SetOutput(logFile)
	}

	logEntry = log.WithFields(log.Fields{
		"app": config.AppName,
	})

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *log.Entry {
	return logEntry.WithField(log.ErrorKey, err)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *log.Entry {
	return logEntry.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields log.Fields) *log.Entry {
	return logEntry.WithFields(fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logEntry.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logEntry.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logEntry.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logEntry.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logEntry.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logEntry.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logEntry.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	logEntry.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logEntry.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logEntry.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logEntry.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logEntry.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logEntry.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logEntry.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logEntry.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	logEntry.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logEntry.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	logEntry.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logEntry.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logEntry.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logEntry.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logEntry.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logEntry.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	logEntry.Fatalln(args...)
}
