package htl

import (
	"os"

	"github.com/Valdenirmezadri/core-go/safe"
	logging "github.com/Valdenirmezadri/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _to safe.Item[*log]

type log struct {
	options Options
	logging logging.Logger
	close   func() error
}

func init() {
	_to = safe.NewItem[*log]()
}

func Start(ops ...Optfunc) (close func() error, err error) {
	o := defaultOps()
	for _, fn := range ops {
		fn(&o)
	}

	return Stop, start(o)
}

func Stop() error {
	close := _to.Get().close
	if close == nil {
		return nil
	}

	return close()
}

func SetLevel(lv string) {
	log := _to.Get()
	options := log.options
	options.level = options.level.New(lv)

	Stop()
	start(options)
}

func Log() logging.Logger {
	return _to.Get().logging
}

func start(o Options) (err error) {
	logger, err := logging.GetLogger(o.module)
	if err != nil {
		return err
	}

	data := &log{logging: logger, options: o}

	if err := data.init(); err != nil {
		return err
	}

	_to.Set(data)

	return nil
}

func (l *log) init() error {
	if l.options.mode == "dev" {
		close, err := l.devLog()
		if err != nil {
			return err
		}

		l.close = close
		return nil
	}

	close, err := l.prodLog()
	if err != nil {
		return err
	}

	l.close = close
	return nil
}

func (l *log) devLog() (close func() error, err error) {
	console := logging.NewLogBackend(os.Stderr, "", 0)
	consoleFormatter := logging.NewBackendFormatter(console, formatConsole)
	consoleBackend := logging.AddModuleLevel(consoleFormatter)
	consoleBackend.SetLevel(l.options.level, l.options.module)

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
	fileBackend.SetLevel(l.options.level, l.options.module)

	return fileBackend, close, nil
}

func (l *log) writerToWithRotation() (writer *logging.LogBackend, close func() error, err error) {
	rotate := &lumberjack.Logger{
		Filename:   l.options.pathLog,
		MaxSize:    int(l.options.maxAge),
		MaxBackups: int(l.options.maxBackups),
		MaxAge:     int(l.options.maxAge),
		Compress:   l.options.compress,
	}

	return logging.NewLogBackend(rotate, "", 0), rotate.Close, nil
}

var formatConsole = logging.MustStringFormatter(
	`%{color} %{time:15:04:05.000} %{shortfile} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var formatFile = logging.MustStringFormatter(
	`%{time:Jan 02 2006 15:04:05} %{shortfile} ▶ %{level:.4s} %{message}`,
)
