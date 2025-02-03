package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
	"resource-update/logger"
	"sync"
)

var (
	configMutex sync.Mutex
	Conf        Config
)

var (
	BootstrapConf        BootstrapConfig
	bootstrapConfigMutex sync.Mutex
)

type Config struct {
	Update []struct {
		GitPath         string `yaml:"git-path"`
		TargetDirectory string `yaml:"target-directory"`
	} `yaml:"update"`
}

type BootstrapConfig struct {
	Network struct {
		Enable   bool   `yaml:"enable"`
		PingAddr string `yaml:"ping-addr"`
	} `yaml:"network"`
	Update struct {
		Resource   string `yaml:"resource"`
		SelfUpdate string `yaml:"self-update"`
	} `yaml:"update"`
}

func parseContent2Config(content []byte) Config {
	var config Config
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		log.Panicln(err)
	}
	return config
}

func parseContent2BootstrapConfig(content []byte) BootstrapConfig {
	var config BootstrapConfig
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		log.Panicln(err)
	}
	return config
}

func loadBootstrapConfig() {
	configFileName := "config.yml"
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		logger.LogErr("配置文件 %s 不存在", configFileName)
		log.Panicln(err)
	}
	content, err := os.ReadFile(configFileName)
	if err != nil {
		log.Panicln(err)
	}
	config := parseContent2BootstrapConfig(content)
	logger.LogInfo("配置文件模式: 本地离线配置文件")
	bootstrapConfigMutex.Lock()
	BootstrapConf = config
	bootstrapConfigMutex.Unlock()
}

func loadUpdateConfig() {
	request, _ := http.NewRequest("GET", BootstrapConf.Update.Resource, nil)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Panicln(err)
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		config := parseContent2Config(bodyBytes)
		configMutex.Lock()
		Conf = config
		configMutex.Unlock()
	} else {
		logger.LogErr("配置文件源错误: %s", BootstrapConf.Update.Resource)
		logger.LogErr("HTTP Error Code: %d", resp.StatusCode)
		log.Panicln(resp.StatusCode)
	}
}

func GetConfig() {
	loadBootstrapConfig()
	loadUpdateConfig()
}
