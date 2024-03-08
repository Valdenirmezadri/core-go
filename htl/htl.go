package htl

import (
	"os"

	logging "github.com/Valdenirmezadri/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var to *log

type log struct {
	Options
	*logging.Logger
	close func() error
}

func Start(ops ...Optfunc) (err error) {
	o := defaultOps()
	for _, fn := range ops {
		fn(&o)
	}

	logger, err := logging.GetLogger(o.module)
	if err != nil {
		return err
	}

	to = &log{
		Logger:  logger,
		Options: o,
	}

	if to.mode == "dev" {
		to.close, err = to.devLog()
		if err != nil {
			return err
		}

		return nil
	}

	to.close, err = to.prodLog()

	return err
}

func Stop() error {
	if to.close == nil {
		return nil
	}

	return to.close()
}

func Log() *logging.Logger {
	return to.Logger
}

func (l *log) devLog() (close func() error, err error) {
	console := logging.NewLogBackend(os.Stderr, "", 0)
	consoleFormatter := logging.NewBackendFormatter(console, formatConsole)
	consoleBackend := logging.AddModuleLevel(consoleFormatter)
	consoleBackend.SetLevel(l.level, l.module)

	fileBackend, close, err := l.fileBackend()
	if err != nil {
		return nil, err
	}

	logging.SetBackend(consoleBackend, fileBackend)
	return close, nil
}

func (l *log) prodLog() (close func() error, err error) {
	fileBackend, close, err := l.fileBackend()
	if err != nil {
		return nil, err
	}

	logging.SetBackend(fileBackend)
	return close, nil
}

func (l *log) fileBackend() (fileFormatter logging.Backend, close func() error, err error) {
	writer, close, err := l.writerToWithRotation()
	if err != nil {
		return nil, nil, err
	}

	fileFromatter := logging.NewBackendFormatter(writer, formatFile)
	fileBackend := logging.AddModuleLevel(fileFromatter)
	fileBackend.SetLevel(l.level, l.module)

	return fileBackend, close, nil
}

func (l *log) writerToWithRotation() (writer *logging.LogBackend, close func() error, err error) {
	rotate := &lumberjack.Logger{
		Filename:   l.pathLog,
		MaxSize:    int(l.maxAge),
		MaxBackups: int(l.maxBackups),
		MaxAge:     int(l.maxAge),
		Compress:   l.compress,
	}

	return logging.NewLogBackend(rotate, "", 0), rotate.Close, nil
}

var formatConsole = logging.MustStringFormatter(
	`%{color} %{time:15:04:05.000} %{shortfile} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var formatFile = logging.MustStringFormatter(
	`%{time:Jan 02 2006 15:04:05} %{shortfile} ▶ %{level:.4s} %{message}`,
)
