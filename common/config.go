package common

import (
	"io/ioutil"
	"log"
	"path"
	"sync"

	"gopkg.in/yaml.v2"
)

// Config config params
type Config struct {
	DB struct {
		DSN    string `yaml:"dsn"`
		Driver string `yaml:"driver"`
	}
}

var configOnce sync.Once

// LoadConfig load configs from yaml file
func LoadConfig() (config Config) {
	configOnce.Do(func() {
		FilePath := GetProjectAbPathByCaller()
		cfgFilePath := path.Join(FilePath, "/common/config.yaml")
		file, err := ioutil.ReadFile(cfgFilePath)
		if err != nil {
			log.Fatalf("ERROR: Could not read config file :%v", err)
		}
		if err := yaml.Unmarshal(file, &config); err != nil {
			// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 	log.Fatal("ERROR: config file not found")
			// 	os.Exit(1)
			// } else {
			log.Fatalf("Fatal error config file: %v", err)
			// }
		}

		return
	})
	return
}
