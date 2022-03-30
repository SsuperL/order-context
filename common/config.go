package common

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// Config config params
type Config struct {
	DB struct {
		DSN    string `yaml:"dsn"`
		Driver string `yaml:"driver"`
	}

	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

var (
	configOnce sync.Once
	// FileConfig instance of file config
	FileConfig Config
)

// LoadConfig load configs from yaml file
func LoadConfig() Config {
	configOnce.Do(func() {
		cfgFilePath := os.Getenv("CONFIG_PATH")
		if _, err := os.Stat(cfgFilePath); errors.Is(err, os.ErrNotExist) {
			log.Fatal("config file does not exist.")
		}
		// FilePath := GetProjectAbPathByCaller()
		// cfgFilePath = path.Join(FilePath, "/common/config.yaml")
		file, err := ioutil.ReadFile(cfgFilePath)
		if err != nil {
			log.Fatalf("ERROR: Could not read config file :%v", err)
		}
		if err := yaml.Unmarshal(file, &FileConfig); err != nil {
			// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 	log.Fatal("ERROR: config file not found")
			// 	os.Exit(1)
			// } else {
			log.Fatalf("Fatal error config file: %v", err)
			// }
		}

		return
	})
	return FileConfig
}
