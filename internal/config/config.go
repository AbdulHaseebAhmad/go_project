package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"Address" env:"Address" env-required:"true"`
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

type Config struct {
	Env          string `yaml:"env" env:"ENV" env-required:"true"` // anotating from yamls file(struct tags)
	Storage_path string `yaml:"storage_path" env:"storage_path" env-required:"true"`
	Database
	HttpServer `yaml:"http_server" env-required:"true"`
}

// definition is done now we write the logic to parse it
func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH") // The Go function os.Getenv("KEY") retrieves the value of an environment variable by its name ("KEY"), and returns it as a string

	if configPath == "" { // if path from env not found use the flag from terminal, using the flag package,from the terminal user will enter the yaml file path
		flags := flag.String("config", "", "path to the config file") // we are setting the flag which shall read for the path to our yaml
		flag.Parse()                                                  // get the config path
		configPath = *flags
		fmt.Println("The path to the yaml holding environment variable:", configPath)

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, error := os.Stat(configPath); os.IsNotExist(error) { // this checks if the path found above is real as if any file exist there, if it does not then throw error
		log.Fatal("Config file does not exist: ", configPath)
	}

	var cfg Config //
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal("Can not read config file: ", err.Error())
	}

	return &cfg
}
