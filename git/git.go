package git

import (
	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"log"
	"os"
	"resource-update/config"
	"resource-update/logger"
	"strings"
)

func removeOldAndClone(path string, url string) {
	logger.LogWarn("正在重新克隆文件...")
	errRemove := os.RemoveAll(path)
	if errRemove != nil {
		log.Panicln(errRemove)
	}
	_, errClone := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if errClone != nil {
		log.Panicln(errClone)
	}
	logger.LogInfo("克隆完成")
}

func updateOrigin(repo *git.Repository, origin string) {
	_, errCreateR := repo.CreateRemote(&gitConfig.RemoteConfig{
		Name: "origin",
		URLs: []string{origin},
	})
	if errCreateR != nil {
		log.Panicln(errCreateR)
	}
	logger.LogInfo("已添加远程源分支: %s", origin)
}

func checkOrigin(repo *git.Repository, origin string) bool {
	logger.LogInfo("正在检查远程源分支...")
	oldOrigin, errRemote := repo.Remote("origin")
	if errRemote != nil {
		log.Panicln(errRemote)
	}
	if oldOrigin.Config().URLs[0] != origin {
		logger.LogWarn("远程源分支 Origin 与新的源远程分支不一致")
		errDelR := repo.DeleteRemote("origin")
		if errDelR != nil {
			log.Panicln(errDelR)
		}
		logger.LogInfo("已删除旧的源远程分支")
		updateOrigin(repo, origin)
		return true
	}
	return false
}

func openRepo(path string) (*git.Repository, error) {
	repo, errOpen := git.PlainOpen(path)
	if errOpen != nil {
		return nil, errOpen
	}
	if repo == nil {
		logger.LogErr("不能正常打开存储库")
		log.Panicln("不能正常打开存储库")
	}
	return repo, nil
}

func update(path string, url string) {
	logger.LogInfo("正在更新 %s (源: %s)", path, url)
	repo, errOpen := openRepo(path)
	if errOpen != nil || checkOrigin(repo, url) {
		removeOldAndClone(path, url)
		errOpen = nil
		repo, errOpen = openRepo(path)
		if errOpen != nil {
			log.Panicln(errOpen)
		}
	}
	worktree, errWorkTree := repo.Worktree()
	if errWorkTree != nil {
		log.Panicln(errWorkTree)
	}
	errPull := worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if errPull != nil {
		if strings.Contains(errPull.Error(), "already up-to-date") {
			logger.LogInfo("已经是最新更新")
			return
		}
		logger.LogErr("本地分支与远程分支产生了冲突，请解决冲突，或者删除目录[%s]", path)
		log.Panicln(errPull)
	}
	logger.LogInfo("%s - 更新完成", path)
}

func DoUpdate() {
	for _, item := range config.Conf.Update {
		update(item.TargetDirectory, item.GitPath)
	}
}
