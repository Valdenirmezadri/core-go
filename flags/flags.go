package flags

import (
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

type Flager interface {
	IsString() bool
	String() string
	IsInt() bool
	Int() int
	IsBool() bool
	Bool() bool
}

func Parse(flags map[string]Flager) map[string]Flager {
	list := make(map[string]Flager)
	for key, _flag := range flags {
		switch flag := _flag.(type) {
		case *Flag:
			switch flag.Kind {
			case stringKind:
				flag.srtValue = kingpin.Flag(key, flag.Help).Short(flag.Short).Default(flag.DefaultVal).String()
				list[key] = flag
			case intKind:
				flag.intValue = kingpin.Flag(key, flag.Help).Short(flag.Short).Default(flag.DefaultVal).Int()
				list[key] = flag
			case boolKind:
				flag.boolValue = kingpin.Flag(key, flag.Help).Short(flag.Short).Default(flag.DefaultVal).Bool()
				list[key] = flag
			}

		}

	}

	kingpin.Parse()

	return list
}

type Kind string

const (
	stringKind Kind = "string"
	intKind    Kind = "int"
	boolKind   Kind = "bool"
)

type Flag struct {
	Kind       Kind
	Help       string
	Short      rune
	DefaultVal string
	srtValue   *string
	intValue   *int
	boolValue  *bool
}

func (f Flag) IsString() bool {
	return f.Kind == stringKind
}

func (f Flag) String() string {
	if f.srtValue != nil {
		return strings.TrimSpace(*f.srtValue)
	}

	return ""
}

func (f Flag) IsInt() bool {
	return f.Kind == intKind
}

func (f Flag) Int() int {
	return *f.intValue
}

func (f Flag) IsBool() bool {
	return f.Kind == boolKind
}
func (f Flag) Bool() bool {
	return *f.boolValue
}

func String(help string, short rune, defaultVal string) Flager {
	return create(stringKind, help, short, defaultVal)
}

func Int(help string, short rune, defaultVal string) Flager {
	return create(intKind, help, short, defaultVal)
}

func Bool(help string, short rune, defaultVal string) Flager {
	return create(boolKind, help, short, defaultVal)
}

func create(kind Kind, help string, short rune, defaultVal string) Flager {
	return &Flag{
		Kind:       kind,
		Help:       help,
		Short:      short,
		DefaultVal: defaultVal,
	}
}
