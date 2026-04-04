package i18n

import (
	"fmt"
	"io/fs"

	"github.com/Valdenirmezadri/core-go/context/corectx"
)

type Language string

const (
	PtBR Language = "pt-BR"
	En   Language = "en"
)

func (l Language) Valid() error {
	if l == PtBR {
		return nil
	}

	if l == En {
		return nil
	}

	return fmt.Errorf("language %s not supported", l)
}

func (l Language) ValidOrDefault() Language {
	if err := l.Valid(); err != nil {
		return PtBR
	}

	return l
}

func (l Language) String() string {
	return string(l)
}

type Switch interface {
	S(ctx corectx.Context, id string, data map[string]string) string
}

type handler struct {
	en   Translator
	ptBR Translator
}

func New(files fs.FS, logger func() Logger) (Switch, error) {
	en, err := newTranslator(files, En.String(), logger)
	if err != nil {
		return nil, err
	}
	ptBR, err := newTranslator(files, PtBR.String(), logger)
	if err != nil {
		return nil, err
	}

	return &handler{en: en, ptBR: ptBR}, nil
}

func (h *handler) S(ctx corectx.Context, id string, data map[string]string) string {
	if ctx != nil {
		lang := Language(ctx.User().Lang)

		if lang == PtBR {
			return h.ptBR.T(id, data)
		}

		if lang == En {
			return h.en.T(id, data)
		}
	}

	return h.en.T(id, data)
}
