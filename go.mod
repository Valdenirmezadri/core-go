module github.com/Valdenirmezadri/core-go

go 1.23

toolchain go1.23.7

require (
	github.com/Valdenirmezadri/viper v1.7.3
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/hashicorp/go-version v1.7.0
	github.com/jackc/pgx/v5 v5.7.4
	github.com/patrickmn/go-cache v2.1.0+incompatible
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.6
)

require (
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/smartystreets/goconvey v1.8.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	go.mau.fi/whatsmeow v0.0.0-20240726213518-bb5852f056ca // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
)

require (
	github.com/Valdenirmezadri/ht-logging v0.0.3
	github.com/fsnotify/fsnotify v1.7.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/labstack/echo/v4 v4.12.0
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.9.2 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace gorm.io/gorm => gorm.io/gorm v1.25.6

replace github.com/Valdenirmezadri/ht-logging => /home/junior/dev/go/ht-logging
