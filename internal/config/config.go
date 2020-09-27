package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port   string `json:"port"`
	AllowedHosts []string `json:"allowed_hosts"`
	
	DBHost string `json:"dbhost"`
	DBPort string `json:"dbport"`
	DBUser string `json:"dbuser"`
	DBPass string `json:"dbpassword"`
	DBName string `json:"dbname"`
	DBMaxConn int `json:"db_max_conn"`
	DBAcquireTimeout int `json:"db_acquire_timeout"`
	
	DocPath string `json:"doc_path"`
	ValidRequestTimeout int `json:"valid_request_timeout"`

}

func NewConfig(pathToConfig string) (*Config, error) {
	conf := new(Config)
	configFile, err := os.Open(pathToConfig)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = configFile.Close()
	}()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
