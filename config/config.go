package repository

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Valdenirmezadri/core-go/observer"
	"github.com/Valdenirmezadri/core-go/safe"
	"github.com/Valdenirmezadri/viper"
	"github.com/fsnotify/fsnotify"
)

type Config any

type Configer[T Config] interface {
	Get() T
	Subscribe(fn func(T)) uint32
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type config[T Config] struct {
	fileReader *viper.Viper
	params     safe.Item[T]
	pub        observer.Publisher[T]
}

func New[T Config](pathFileName string) (Configer[T], error) {

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
		params:     safe.NewItem[T](),
		pub:        observer.NewPublisher[T](true),
	}

	err = repo.init()

	if err != nil {
		return nil, err
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		err := repo.init()
		if err != nil {
			log.Default().Printf("error on read file config %+v", err)
			return
		}

		repo.pub.Next(repo.Get())
	})

	return repo, nil
}

func (r config[T]) init() error {
	var c T

	err := r.fileReader.Unmarshal(&c)
	if err != nil {
		return err
	}

	r.params.Set(c)
	return nil
}

func (r *config[T]) Subscribe(fn func(T)) uint32 {
	return r.pub.Subscribe(observer.NewListener[T](fn))
}

func (r config[T]) Get() T {
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
