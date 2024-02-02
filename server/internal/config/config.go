package config

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Args struct {
	ConfigPath string
}

type Config struct {
	App struct {
		Url string `yaml:"url" env:"APP_URL"`
	} `yaml:"app"`
	Telegram struct {
		Token string `yaml:"token" env:"TELEGRAM_TOKEN"`
	} `yaml:"telegram"`
	HTTP struct {
		IP   string `yaml:"ip" env:"HTTP_IP"`
		Port string `yaml:"port" env:"HTTP_PORT"`
		CORS struct {
			AllowedMethods     []string `yaml:"allowed-methods" env:"HTTP_CORS_ALLOWED_METHODS"`
			AllowedOrigins     []string `yaml:"allowed-origins" env:"HTTP_CORS_ALLOWED_ORIGINS"`
			AllowCredentials   bool     `yaml:"allowed-credentials" env:"HTTP_CORS_ALLOWED_CREDENTIALS"`
			AllowedHeaders     []string `yaml:"allowed-headers" env:"HTTP_CORS_ALLOWED_HEADERS"`
			OptionsPassthrough bool     `yaml:"options-passthrough" env:"HTTP_CORS_OPTIONS_PASSTHROUGH"`
			ExposedHeaders     []string `yaml:"exposed-headers" env:"HTTP_CORS_EXPOSED_HEADERS"`
			Debug              bool     `yaml:"debug" env:"HTTP_CORS_DEBUG"`
		} `yaml:"cors"`
	} `yaml:"http"`
	PostgresSQL struct {
		Username string `yaml:"username" env:"PSQL_USERNAME" env-default:"postgres"`
		Password string `yaml:"password" env:"PSQL_PASSWORD" env-default:"postgres"`
		Host     string `yaml:"host" env:"PSQL_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"PSQL_PORT" env-default:"5432"`
		Database string `yaml:"database" env:"PSQL_DATABASE" env-default:"postgres"`
	} `yaml:"postgres"`
}

var instance Config
var once sync.Once

func GetConfig() *Config {
	log.Println("config init")
	once.Do(func() {
		args := processArgs(&instance)
		if err := cleanenv.ReadConfig(args.ConfigPath, &instance); err != nil {
			log.Println(err)
			os.Exit(2)
		}
	})
	return &instance
}

func processArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", "configs/config.yml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		log.Println(f.Output())
		log.Println(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}
