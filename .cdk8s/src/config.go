package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"Environment"`
	Outdir      string `mapstructure:"Outdir"`
	StorageName string `mapstructure:"StorageName"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	viper.SetConfigFile("config.json")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("CDK8S")
	viper.AutomaticEnv()
	viper.BindEnv("Environment")

	pflag.StringP("environment", "e", viper.GetString("Environment"), "Установка окружения (например, Development, Production)")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Чтение конфигурации
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	if environment := viper.GetString("environment"); len(environment) != 0 {
		// Переопределение базовой конфигурации файлом окружения
		viper.SetConfigName("config." + environment)

		viper.MergeInConfig()
	}

	// Распаковка конфигурации в структуры
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("ошибка распаковки конфигурации: %w", err)
	}

	return config, nil
}
