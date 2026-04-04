package i18n

import (
	"io/fs"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Logger interface {
	Warningf(format string, args ...interface{})
}

type Translator interface {
	T(id string, data map[string]string) string
}

type Bundle struct {
	bundle *i18n.Bundle
	loc    *i18n.Localizer
	logger func() Logger
}

func newTranslator(files fs.FS, lang string, logger func() Logger) (Translator, error) {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		b.LoadMessageFileFS(files, path)
		return nil
	})

	t := &Bundle{
		bundle: b,
		loc:    i18n.NewLocalizer(b, lang),
		logger: logger,
	}

	return t, nil
}

func (b *Bundle) T(id string, data map[string]string) string {
	msg, err := b.loc.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil && b.logger != nil {
		b.logger().Warningf("i18n: translation missing for '%s': %v", id, err)
		return id
	}

	return msg
}
