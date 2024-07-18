package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var Config config

type config struct {
	BasePath  string
	Logger    Logger
	Database  map[string]Database
	Listen    string
	ServerCrt    string
	ServerKey    string
	APISecret string
	AppEnv    string
	AppName   string
	Redis     map[string]RedisConfig
}

type Logger struct {
	Debug        bool
	OutFile      bool
	LogPath      string
	MaxAge       time.Duration //日志最大保存时间单位小时
	RotationTime time.Duration //日志切割时间间隔单位小时
}

type Database struct {
	DriverName     string
	DataSourceName string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (c *config) Init() {
	var appenv string
	if len(os.Args) < 2 {
		appenv = os.Getenv("APP_ENV")
	} else {
		appenv = os.Args[1]
	}

	if appenv == "" {
		Config.AppEnv = "dev"
	} else {
		Config.AppEnv = appenv
	}
	exPath,_ := os.Getwd()
	fmt.Println("expath",exPath)
	separatorString := string(os.PathSeparator)
	// /src/src/
	projectDir:=exPath
	crtBaseUrl:=projectDir+separatorString+"src"+separatorString
	configPatchStr:=projectDir+separatorString+"src"+separatorString+"config." + Config.AppEnv + ".json"
	if strings.Contains(exPath,"go-websocket-broadcast"+separatorString+"src"){
		configPatchStr=projectDir+separatorString+"config." + Config.AppEnv + ".json"
		crtBaseUrl=projectDir+separatorString
	}


	log.Info("configPatchStr=",configPatchStr)
	file, err := os.Open(configPatchStr)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	c.ServerCrt=crtBaseUrl+c.ServerCrt
	c.ServerKey=crtBaseUrl+c.ServerKey
	cwd, _ := os.Getwd()
	if Config.BasePath != "" {
		cwd = Config.BasePath
	} else {
		Config.BasePath = cwd
	}
}
