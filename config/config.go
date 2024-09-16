package config

import (
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/memberlist"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func GetNodeName() string {
	if name := viper.GetString("node.name"); name != "" {
		return name
	}

	return uuid.NewString()
}

func GetNodeConfig() *memberlist.Config {
	config := viper.GetString("node.config")

	switch config {
	case "lan":
		return memberlist.DefaultLANConfig()
	case "wan":
		return memberlist.DefaultWANConfig()
	default:
		return memberlist.DefaultLocalConfig()
	}
}

func GetLogLevel() *logutils.LevelFilter {
	level := "WARN"

	if viper.GetString("log.level") != "" {
		level = viper.GetString("log.level")
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(level),
		Writer:   os.Stderr,
	}

	return filter
}

func Init(flag *pflag.FlagSet) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "."))

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.For("config", "init").Debug("No config file found, use from environment variables")
		} else {
			log.For("config", "init").Fatal(err)
		}
	}

	bindFlags(flag)
}

func bindFlags(flag *pflag.FlagSet) {
	flag.VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", ".")

		if f.Changed && !viper.IsSet(configName) {
			viper.BindPFlag(configName, f)
		}
	})
}
