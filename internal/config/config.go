package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xiaolingzi/lingorm/internal/common"
)

// DatabaseInfoEntity struct
type DatabaseInfoEntity struct {
	Driver   string
	Host     string
	Port     string
	Database string
	User     string
	Password string
	Charset  string
	Servers  []serverEntity
}

type serverEntity struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	Charset  string
	Mode     string
	RWeight  int
	WWeight  int
	Weight   int
}

var databaseConfig map[string]DatabaseInfoEntity

// GetDatabaseInfo returns the database infomation
func GetDatabaseInfo(key string, mode string) DatabaseInfoEntity {
	databaseInfo := getDatabaseInfoByKey(key)
	serverLength := len(databaseInfo.Servers)
	if serverLength == 0 || (mode != common.DbWriteMode && mode != common.DbReadMode) {
		return databaseInfo
	}

	var serverList []serverEntity

	for i := 0; i < serverLength; i++ {
		server := databaseInfo.Servers[i]
		if mode == common.DbWriteMode {
			server.Weight = server.WWeight
		} else if mode == common.DbWriteMode {
			server.Weight = server.RWeight
		}
		if server.Weight <= 0 {
			continue
		}

		if strings.Contains(server.Mode, mode) {
			serverList = append(serverList, server)
		} else if server.Mode == "" && mode == common.DbReadMode {
			serverList = append(serverList, server)
		}
	}

	var targetServer serverEntity
	if len(serverList) == 0 {
		return databaseInfo
	}
	if len(serverList) == 1 {
		targetServer = serverList[0]
	} else {
		targetServer = getRandomDatabase(serverList)
	}

	if len(targetServer.Host) > 0 {
		databaseInfo.Host = targetServer.Host
	}
	if len(targetServer.Port) > 0 {
		databaseInfo.Port = targetServer.Port
	}
	if len(targetServer.Database) > 0 {
		databaseInfo.Database = targetServer.Database
	}
	if len(targetServer.User) > 0 {
		databaseInfo.User = targetServer.User
	}
	if len(targetServer.Password) > 0 {
		databaseInfo.Password = targetServer.Password
	}
	if len(targetServer.Charset) > 0 {
		databaseInfo.Charset = targetServer.Charset
	}

	return databaseInfo
}

// GetDatabaseDriver returns the database driver
func GetDatabaseDriver(key string) string {
	databaseInfo := getDatabaseInfoByKey(key)
	return databaseInfo.Driver
}

func getRandomDatabase(serverList []serverEntity) serverEntity {
	var result serverEntity
	if len(serverList) == 0 {
		return result
	}

	if len(serverList) == 1 {
		return serverList[0]
	}

	sum := 0
	for _, server := range serverList {
		if server.Weight <= 0 {
			continue
		}
		sum += server.Weight
		server.Weight = sum
	}

	rand.Seed(time.Now().Unix())
	index := 0
	if sum == 0 {
		index = rand.Intn(len(serverList))
	} else {
		ranNum := rand.Intn(sum) + 1
		for i := 0; i < len(serverList); i++ {
			if ranNum <= serverList[i].Weight {
				index = i
			}
		}
	}
	result = serverList[index]
	return result
}

func getDatabaseConfig() map[string]DatabaseInfoEntity {
	if databaseConfig != nil {
		return databaseConfig
	}
	filename := getConfigFilename()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	databaseConfig = make(map[string]DatabaseInfoEntity)
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

func getConfigFilename() string {
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
