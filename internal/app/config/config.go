package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/smf8/arvan-voucher/pkg/database"
	"github.com/smf8/arvan-voucher/pkg/router"

	"github.com/sirupsen/logrus"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

const _Prefix = "WALLET_"

type Config struct {
	LogLevel string                  `koanf:"log_level"`
	Server   router.ServerConfig     `koanf:"server"`
	Database database.DatabaseConfig `koanf:"database"`
}

var def = Config{
	LogLevel: "debug",
	Server: router.ServerConfig{
		Port:         ":8001",
		Debug:        true,
		NameSpace:    "wallet",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	},

	Database: database.DatabaseConfig{
		ConnectionAddress:  "postgresql://root@127.0.0.1:26257/defaultdb",
		RetryDelay:         time.Second,
		MaxRetry:           20,
		ConnectionLifetime: 30 * time.Minute,
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
		LogLevel:           4,
	},
}

func New() Config {
	var instance Config

	k := koanf.New(".")

	if err := k.Load(structs.Provider(def, "koanf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		logrus.Errorf("error loading file: %s", err)
	}

	if err := k.Load(env.Provider(_Prefix, ".", func(s string) string {
		parsedEnv := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, _Prefix)), "__", "-")

		fmt.Println(parsedEnv)

		return strings.ReplaceAll(parsedEnv, "_", ".")
	}), nil); err != nil {
		logrus.Errorf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	logrus.Infof("following configuration is loaded:\n%+v", instance)

	return instance
}
