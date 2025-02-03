package main

import (
	"bufio"
	"fmt"
	"os"
	"resource-update/config"
	"resource-update/git"
	"resource-update/network"
	"resource-update/version"
)

func main() {
	defer func() {
		_ = recover()
		fmt.Print("按回车（Enter）结束...")
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	}()
	config.GetConfig()
	network.CheckNetwork()
	version.CheckUpdate()
	git.DoUpdate()
}
