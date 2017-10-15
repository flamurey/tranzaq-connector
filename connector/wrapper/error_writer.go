package wrapper

import (
	"github.com/flamurey/tranzaq-connector/connector/logger"
	"github.com/sirupsen/logrus"
)

type ErrorWriter struct {
	logger *logrus.Entry
}

func CreateWrapperOut() ErrorWriter {
	return ErrorWriter{
		logger: logger.WithField("app", "tranzaq-wrapper"),
	}
}

func (ew ErrorWriter) Write(p []byte) (n int, err error) {
	ew.logger.Error(string(p))
	return len(p), nil
}
