package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" required:"true"`
	GRPC       GRPC       `yaml:"grpc" required:"true"`
	RedisDB    Redis      `yaml:"redis"`
	DB         DataBase   `yaml:"database"`
	Prometheus Prometheus `yaml:"prometheus"`
}
type Prometheus struct {
	Timeout     time.Duration `yaml:"timeout"  env-default:"4s"`
	Host        string        `yaml:"host"  env-default:"localhost"`
	Port        string        `yaml:"port"  env-default:"8081"`
	Debug       bool          `yaml:"debug"  env-default:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"60s"`
}
type GRPC struct {
	Port    string        `yaml:"port" required:"true"`
	Host    string        `yaml:"host" required:"true"`
	Timeout time.Duration `yaml:"timeout" required:"true"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Retries  int    `yaml:"retries"`
	DBNumber int    `yaml:"dbnumber"`
}
type DataBase struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Database    string `yaml:"dbname"`
	SSL         string `yaml:"ssl"`
	MaxAttempts int    `yaml:"max_attempts"`
}

func InitConfig() *Config {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env.dev"
	}
	fmt.Println("env name", envFile)
	if err := godotenv.Load(envFile); err != nil {
		slog.Error("ошибка при инициализации переменных окружения", err.Error())
	}
	configPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist:%s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
