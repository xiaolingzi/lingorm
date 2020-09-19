package sqlite

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/xiaolingzi/lingorm/internal/common"
	"github.com/xiaolingzi/lingorm/internal/config"
)

// DatabaseInfoEntity struct
type Config struct {
	Driver   string
	File     string
	User     string
	Password string
	Crypt    string
	Salt     string
	Timeout  int
}

func NewConfig() *Config {
	return new(Config)
}

var databaseConfig map[string]Config

// GetDatabaseInfo returns the database infomation
func (c *Config) GetDatabaseInfo(key string) Config {
	databaseInfo := c.getDatabaseInfoByKey(key)
	return databaseInfo
}

func (c *Config) getDatabaseConfig() map[string]interface{} {
	filename := config.GetConfigFilename()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	databaseConfig := make(map[string]interface{})
	json.Unmarshal([]byte(content), &databaseConfig)
	return databaseConfig
}

func (c *Config) getDatabaseInfoByKey(key string) Config {
	if databaseConfig != nil {
		if _, ok := databaseConfig[key]; ok {
			return databaseConfig[key]
		}
	}
	configMap := c.getDatabaseConfig()
	var result Config
	if _, ok := configMap[key]; ok {
		tempValue, _ := json.Marshal(configMap[key])
		json.Unmarshal(tempValue, &result)
	}
	if result.File == "" {
		common.NewError().Throw("Database config not found or invalid.")
	}
	return result
}
