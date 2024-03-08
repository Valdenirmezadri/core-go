package htl

import "github.com/Valdenirmezadri/go-logging"

type Optfunc func(*Options)

type Options struct {
	module     string
	level      logging.Level
	mode       string
	pathLog    string
	maxSize    uint
	maxBackups uint
	maxAge     uint
	compress   bool
}

func Level(l string) Optfunc {
	return func(o *Options) {
		o.level = o.level.New(l)
	}
}

func AsProd(o *Options) {
	o.mode = "prod"
}

func Module(module string) Optfunc {
	return func(o *Options) {
		o.module = module
	}
}

func PathLog(pathFileName string) Optfunc {
	return func(o *Options) {
		o.pathLog = pathFileName
	}
}

/*
MaxSize is the maximum size in megabytes of the log file before it gets rotated.
It defaults to 100 megabytes
*/
func MaxSize(n uint) Optfunc {
	return func(o *Options) {
		o.maxSize = n
	}
}

/*
MaxBackups is the maximum number of old log files to retain.
The default is to retain 10 (though MaxAge may still cause them to get deleted.)
*/
func MaxBackups(n uint) Optfunc {
	return func(o *Options) {
		o.maxBackups = n
	}
}

/*
MaxAge is the maximum number of days to retain old log files
based on the timestamp encoded in their filename.
Note that a day is defined as 24 hours and may not exactly correspond
to calendar days due to daylight savings, leap seconds, etc.
The default is 365 days.
*/
func MaxAge(n uint) Optfunc {
	return func(o *Options) {
		o.maxAge = n
	}
}

/*
Compress determines if the rotated log files should be compressed using gzip.
The default is not to perform compression.
*/
func WithCompress() Optfunc {
	return func(o *Options) {
		o.compress = true
	}
}

func defaultOps() Options {
	return Options{
		module:     "htlog",
		mode:       "dev",
		pathLog:    "htlog.log",
		level:      logging.DEBUG,
		maxSize:    100,
		maxBackups: 10,
		maxAge:     365,
		compress:   false,
	}
}
