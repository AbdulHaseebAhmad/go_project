package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  bool
}
type Config struct {
	Env          string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"` // anotating from yamls file(struct tags)
	Storage_path string `yaml:"storage_path" env:"storage_path" env-required:"true" env-default:"production"`
	Database     Database
	HttpServer   `yaml:"http_server" env-required:"true"`
}

// definition is done no we write the logic to parse it
func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH") // get the path from env

	if configPath == "" { // if path from env not found use the flag from terminal using the flag package, remember at job to start the app we were setting everything from terminal for the react app
		flags := flag.String("config", "", "path to the config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, error := os.Stat(configPath); os.IsNotExist(error) {
		log.Fatal("Config file does not exist: ", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal("Can not read config file: ", err.Error())
	}

	return &cfg
}
