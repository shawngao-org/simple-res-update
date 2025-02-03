package network

import (
	"log"
	"os/exec"
	"resource-update/config"
	"resource-update/logger"
)

func CheckNetwork() {
	if !config.BootstrapConf.Network.Enable {
		logger.LogWarn("跳过了网络检查.")
		return
	}
	logger.LogWarn("正在检查网络，请稍候...")
	_, err := exec.Command("ping", "-c", "4", config.BootstrapConf.Network.PingAddr).CombinedOutput()
	if err != nil {
		logger.LogErr("不能连接到互联网，请检查网络或配置文件的检测地址是否正确")
		log.Panicln(err)
	}
	logger.LogInfo("网络连接成功")
}
