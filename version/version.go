package version

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"resource-update/config"
	"resource-update/logger"
	"strconv"
	"strings"
	"sync"
)

const VERSION = "1.0.0"
const TmpPath = "update-windows-amd64-tmp.exe"
const UpdateHelperPath = "update-helper.bat"

var (
	VersionConf        VersionConfig
	versionConfigMutex sync.Mutex
)

type VersionConfig struct {
	Version string `yaml:"version"`
	Update  string `yaml:"update"`
}

func parseContent2VersionConfig(content []byte) VersionConfig {
	var conf VersionConfig
	err := yaml.Unmarshal(content, &conf)
	if err != nil {
		log.Panicln(err)
	}
	return conf
}

func loadVersionConfig() {
	request, _ := http.NewRequest("GET", config.BootstrapConf.Update.SelfUpdate, nil)
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
		conf := parseContent2VersionConfig(bodyBytes)
		versionConfigMutex.Lock()
		VersionConf = conf
		versionConfigMutex.Unlock()
	} else {
		logger.LogErr("配置文件源错误: %s", config.BootstrapConf.Update.SelfUpdate)
		logger.LogErr("HTTP Error Code: %d", resp.StatusCode)
		log.Panicln(resp.StatusCode)
	}
}

func compareVersion(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		num1, _ := strconv.Atoi(parts1[i])
		num2, _ := strconv.Atoi(parts2[i])
		if num1 < num2 {
			return -1 // v1 < v2
		} else if num1 > num2 {
			return 1 // v1 > v2
		}
	}
	if len(parts1) < len(parts2) {
		return -1 // v1 < v2 (1.2 < 1.2.1)
	} else if len(parts1) > len(parts2) {
		return 1 // v1 > v2 (1.2.1 > 1.2)
	}
	return 0 // v1 == v2
}

func CheckUpdate() {
	logger.LogInfo("当前版本: %s", VERSION)
	logger.LogInfo("正在检查更新...")
	loadVersionConfig()
	res := compareVersion(VERSION, VersionConf.Version)
	if res != 0 {
		logger.LogInfo("发现新版本: %s", VersionConf.Version)
		downloadUpgradePackage()
		doUpgrade()
	}
	logger.LogInfo("已是最新版: %s", VERSION)
}

func downloadUpgradePackage() {
	logger.LogInfo("开始下载更新包，如果下载失败，可以前往该连接下载: %s", VersionConf.Update)
	request, _ := http.NewRequest("GET", VersionConf.Update, nil)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Panicln(err)
	}
	if resp.StatusCode == http.StatusOK {
		file, err := os.Create(TmpPath)
		if err != nil {
			log.Panicln(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Panicln(err)
			}
		}(file)
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		logger.LogErr("配置文件源错误: %s", config.BootstrapConf.Update.SelfUpdate)
		logger.LogErr("HTTP Error Code: %d", resp.StatusCode)
		log.Panicln(resp.StatusCode)
	}
}

func doUpgrade() {
	logger.LogInfo("开始重启升级...")
	cmd := exec.Command("cmd", "/C", UpdateHelperPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
