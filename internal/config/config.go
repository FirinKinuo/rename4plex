package config

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Debug               bool     `yaml:"debug"`
	LogLevel            string   `yaml:"log_level"`
	LogPath             string   `yaml:"log_path"`
	DirPlexAnimeLibrary string   `yaml:"dir_plex_anime_library"`
	RegexpsAnimeData    []string `yaml:"regexps_anime_data"`
}

func GetConfig() *Config {
	configPath := filepath.FromSlash("/etc/go-plex-anime/config.yaml")
	var config *Config
	var once sync.Once

	once.Do(func() {
		config = &Config{}
		if err := cleanenv.ReadConfig(configPath, config); err != nil {
			logrus.Errorf("ERROR || Unable to initialize configuration file: %s!Reason: %s", configPath, err)
			os.Exit(-1)
		}
	})
	return config
}

func InitLogger() {
	cfg := GetConfig()
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)

	if !cfg.Debug {
		logFilePath := filepath.FromSlash(cfg.LogPath)

		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			logrus.Errorf("Unable to use log file %s! Reason:%s", cfg.LogPath, err)
		}
		logrus.SetOutput(io.MultiWriter(os.Stderr, logFile))
	}

	if err != nil {
		logrus.Errorf("Unable to recognize logging level %s! INFO level will be used", cfg.LogLevel)
		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)
}
