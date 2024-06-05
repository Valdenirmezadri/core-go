package repository

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Valdenirmezadri/core-go/htl"
	"github.com/Valdenirmezadri/core-go/observer"
	"github.com/Valdenirmezadri/core-go/safe"
	"github.com/Valdenirmezadri/viper"
	"github.com/fsnotify/fsnotify"
)

type Config any

type Configer[T Config] interface {
	GetCold() T
	Get() T
	Subscribe(fn func(T)) uint32
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type config[T Config] struct {
	fileReader *viper.Viper
	first      safe.Item[bool]
	coldParams safe.Item[T]
	params     safe.Item[T]
	defaultsFn setDefaults[T]
	pub        observer.Publisher[T]
}

type setDefaults[T Config] func(c T) (T, error)

func WithDefaults[T Config](pathFileName string, fn setDefaults[T]) (Configer[T], error) {
	config, err := new[T](pathFileName)
	if err != nil {
		return nil, err
	}

	config.defaultsFn = fn

	return config, config.beforeReturn()
}

func New[T Config](pathFileName string) (Configer[T], error) {
	config, err := new[T](pathFileName)
	if err != nil {
		return nil, err
	}

	return config, config.beforeReturn()
}

func new[T Config](pathFileName string) (*config[T], error) {

	directory, name, extension, err := getFile(pathFileName)
	if err != nil {
		return nil, err
	}

	viper := viper.New()

	viper.SetConfigType(extension)
	viper.AddConfigPath(directory)
	viper.SetConfigName(name)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err

	}

	viper.WatchConfig()

	repo := &config[T]{
		fileReader: viper,
		first:      safe.NewItemWithData(false),
		coldParams: safe.NewItem[T](),
		params:     safe.NewItem[T](),
		pub:        observer.NewPublisher[T](true),
	}

	return repo, nil
}

func (r *config[T]) beforeReturn() error {
	if err := r.init(); err != nil {
		return err
	}

	r.fileReader.OnConfigChange(func(in fsnotify.Event) {
		err := r.init()
		if err != nil {
			log.Default().Printf("error on read file config %w", err)
			return
		}

		r.pub.Next(r.Get())
	})

	return nil
}

func (r *config[T]) init() error {
	var c T

	err := r.fileReader.Unmarshal(&c)
	if err != nil {
		return err
	}

	if r.defaultsFn != nil {
		c, err = r.defaultsFn(c)
		if err != nil {
			if htl.Log() != nil {
				htl.Log().Error(err)
				return err
			}

			log.Default().Println(err)
			return err
		}
	}

	if !r.first.Get() {
		r.first.Set(true)
		r.coldParams.Set(c)
	}

	r.params.Set(c)

	return nil
}

func (r *config[T]) Subscribe(fn func(T)) uint32 {
	return r.pub.Subscribe(observer.NewListener[T](fn))
}

func (r *config[T]) GetCold() T {
	return r.params.Get()
}

func (r *config[T]) Get() T {
	return r.params.Get()
}

func getExt(pathFileName string) (string, error) {
	ext := filepath.Ext(pathFileName)
	if ext == "" {
		return "", fmt.Errorf("arquivo de configuração %s não possui a extensão yaml", pathFileName)
	}

	ext = ext[1:]
	if ext != "yaml" {
		return "", fmt.Errorf("arquivo de configuração %s não possui a extensão yaml", pathFileName)
	}

	return ext, nil
}

func getfileName(pathFileName string) (string, error) {
	fileName := filepath.Base(pathFileName)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	if fileNameWithoutExt == "" {
		return "", fmt.Errorf("arquivo de configuração %s não encontrado", pathFileName)
	}

	return fileNameWithoutExt, nil
}

func getFile(pathFileName string) (directoryPath, name, extension string, err error) {
	if _, err = os.Stat(pathFileName); err != nil {
		return "", "", "", err
	}

	directoryPath = filepath.Dir(pathFileName)
	if directoryPath == "" {
		return "", "", "", fmt.Errorf("arquivo de configuração %s não encontrado", pathFileName)
	}

	name, err = getfileName(pathFileName)
	if err != nil {
		return "", "", "", err
	}

	extension, err = getExt(pathFileName)
	if err != nil {
		return "", "", "", err
	}

	return directoryPath, name, extension, nil
}
