package htl

import (
	"fmt"
	"os"

	"github.com/Valdenirmezadri/core-go/safe"
	logging "github.com/Valdenirmezadri/ht-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

const defaultInstance = "defaultInstance"

var _to safe.Lister[string, *log]

type log struct {
	options Options
	logging logging.Logger
	close   func() error
}

func init() {
	_to = safe.NewList[string, *log]()
}

func Start(ops ...Optfunc) (close func() error, err error) {
	return startDefault(ops...)
}

func startDefault(ops ...Optfunc) (close func() error, err error) {
	o := defaultOps()
	for _, fn := range ops {
		fn(&o)
	}

	return Stop, _startDefault(o)
}

func Stop() error {
	var errs []error
	if err := stop(defaultInstance); err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	var allErr error
	for _, err := range errs {
		allErr = fmt.Errorf("%w\n%w", allErr, err)
	}

	return allErr
}

func stop(instance string) error {
	ok, log := _to.Load(instance)
	if !ok {
		return fmt.Errorf("instance log %s not found", instance)
	}

	close := log.close
	if close == nil {
		return nil
	}

	return close()
}

func SetLevel(lv string) {
	_to.Range(func(instance string, _ *log) bool {
		setLevel(instance, lv)
		return true
	})
}

func setLevel(instance, lv string) {
	ok, log := _to.Load(instance)
	if !ok {
		return
	}

	options := log.options
	options.level = options.level.New(lv)

	log.options = options
	log.logging.SetLevel(log.options.level.String())
}

func Log() logging.Logger {
	_, log := _to.Load(defaultInstance)
	return log.logging
}

func _startDefault(o Options) (err error) {
	log, err := _initDefault(o)
	if err != nil {
		return err
	}

	_to.Add(string(defaultInstance), log)

	return nil
}

func _initDefault(o Options) (*log, error) {
	var backEnds []logging.Backend

	file, close, err := prodLog(o)
	if err != nil {
		return nil, err
	}

	backEnds = append(backEnds, file)

	if o.mode == "dev" {
		console, err := devLog()
		if err != nil {
			return nil, err
		}

		backEnds = append(backEnds, console)
	}

	logging, err := logging.New(o.level.String(), backEnds...)
	if err != nil {
		return nil, err
	}

	return &log{logging: logging, close: close, options: o}, nil
}

func devLog() (consoleBackend logging.Backend, err error) {
	console := logging.NewLogBackend(os.Stderr, "", 0)
	consoleFormatter := logging.NewBackendFormatter(console, formatConsole)

	return consoleFormatter, nil
}

func prodLog(o Options) (backend logging.Backend, close func() error, err error) {
	fileBackend, close, err := fileBackend(o)
	if err != nil {
		return nil, nil, err
	}

	return fileBackend, close, nil
}

func fileBackend(o Options) (fileFormatter logging.Backend, close func() error, err error) {
	writer, close, err := writerToWithRotation(o)
	if err != nil {
		return nil, nil, err
	}

	fileFromatter := logging.NewBackendFormatter(writer, formatFile)

	return fileFromatter, close, nil
}

func writerToWithRotation(o Options) (writer *logging.LogBackend, close func() error, err error) {
	rotate := &lumberjack.Logger{
		Filename:   o.pathLog,
		MaxSize:    int(o.maxAge),
		MaxBackups: int(o.maxBackups),
		MaxAge:     int(o.maxAge),
		Compress:   o.compress,
	}

	return logging.NewLogBackend(rotate, "", 0), rotate.Close, nil
}

var formatConsole = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfile} ▶ %{level:.4s} %{color:reset}%{message}`,
)

var formatFile = logging.MustStringFormatter(
	`%{time:Jan 02 2006 15:04:05} %{shortfile} ▶ %{level:.4s} %{message}`,
)
