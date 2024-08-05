package htl

import logging "github.com/Valdenirmezadri/ht-logging"

func StartInstance(instance string, ops ...Optfunc) (err error) {
	o := defaultOps()
	for _, fn := range ops {
		fn(&o)
	}

	log, err := _initInstance(o)
	if err != nil {
		return err
	}

	_to.Add(string(instance), log)

	return nil
}

func _initInstance(o Options) (*log, error) {
	file, close, err := prodLog(o)
	if err != nil {
		return nil, err
	}

	logging, err := logging.New(o.level.String(), file)
	if err != nil {
		return nil, err
	}

	return &log{logging: logging, close: close, options: o}, nil
}

func StopInstance(instance string) error {
	return stop(instance)
}

func In(instance string) logging.Logger {
	if _to.Has(instance) {
		_, log := _to.Load(instance)
		return log.logging
	}

	return Log()
}

func StopInstances() (errors []error) {
	_to.Range(func(key string, value *log) bool {
		if key != defaultInstance {
			errors = append(errors, value.close())
		}

		return true
	})

	return errors
}
