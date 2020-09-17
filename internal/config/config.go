package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiaolingzi/lingorm/internal/common"
)

// DatabaseInfoEntity struct
type DatabaseInfoEntity struct {
	Driver string
}

var drivers map[string]string

// GetDatabaseDriver returns the database driver
func GetDatabaseDriver(key string) string {
	if drivers != nil {
		if driver, ok := drivers[key]; ok {
			return driver
		}
	} else {
		drivers = make(map[string]string)
	}
	databaseInfo := getDatabaseInfoByKey(key)
	drivers[key] = databaseInfo.Driver
	return drivers[key]
}

func getDatabaseConfig() map[string]DatabaseInfoEntity {
	filename := GetConfigFilename()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	databaseConfig := make(map[string]DatabaseInfoEntity)
	json.Unmarshal([]byte(content), &databaseConfig)
	return databaseConfig
}

func getDatabaseInfoByKey(key string) DatabaseInfoEntity {
	configMap := getDatabaseConfig()
	var result DatabaseInfoEntity
	if _, ok := configMap[key]; ok {
		result = configMap[key]
	}
	return result
}

func GetConfigFilename() string {
	appDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	configFilename := os.Getenv("LINGORM_CONFIG")
	if configFilename == "" {
		common.NewError().Throw("Please set env variable LINGORM_CONFIG first before using lingorm.")
	}

	if !strings.HasPrefix(configFilename, "/") && !strings.HasPrefix(configFilename, "\\") && !strings.Contains(configFilename, ":") {
		configFilename = filepath.Join(appDir, configFilename)
	}

	_, err := os.Stat(configFilename)
	if os.IsNotExist(err) {
		common.NewError().Throw("Database config file for lingorm not found!\n" + configFilename)
	}
	return configFilename
}
