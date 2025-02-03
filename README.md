# 这是一个简易的更新程序

## 简介

在玩我的世界服务器的时候，会经常更新一些模组之类的文件（手动替换还是有点麻烦），
由于经常能用到Git这类工具，由此启发，可以用Git进行版本管理。（但是只考虑了Windows系统）

## 功能

+ 可以在线更新资源文件
+ 进行差异化更新，节省流量消耗以及等待时间
+ 有一定的兼容性，任何Git平台均可（目前在Gitea上面使用过）（可以是自建Git，也可以是公有Git）
+ 更新资源的配置全部在Git上面进行管理
+ 这个程序本身也可以进行自我更新（在启动的时候）（同样，也是基于Git的）

## 目录结构

```text
. 项目根目录
├── README.md 自述文件
├── build.sh 构建脚本
├── config 配置
│   └── config.go
├── config.yml 配置文件
├── git Git操作
│   └── git.go
├── go.mod
├── go.sum
├── logger 日志
│   └── logger.go
├── main.go 程序入口点
├── network 网络检查
│   └── check.go
├── update-helper.bat 更新辅助脚本
└── version 自我版本更新
    └── version.go
```

## 使用方法

1. 找个Git，新建2个仓库，一个存放资源（update-mods），一个存放自我更新和配置（update-config）
2. update-mods放资源就可以，没什么说的，主要是update-config
3. 在配置仓库中，新建2个文件

  config.yml
```yaml
# 更新的相关配置
update:
  # 更新资源的仓库地址
  - git-path: https://xxxxxx/ZetoHkr/update-mods.git
    # 更新的目标路径
    target-directory: .minecraft/versions/1.20.1-Forge/mods

```

update.yml
```yaml
# 版本号，用于检测更新
version: 1.0.0
# 更新下载的URL
update: https://xxxxxx/ZetoHkr/update-config/raw/branch/main/resource-update.exe

```

4. 远程仓库配置完成之后，还需要配置本地
5. 本地也有一个配置文件，和该程序放在同一个目录即可，配置内容如下：

config.yml
```yaml
# 网络检查
network:
  # 是否启用
  enable: false
  # ping的IP或者域名
  ping-addr: www.baidu.com
update:
  # 更新源
  resource: https://xxxxxx/ZetoHkr/update-config/raw/branch/main/config.yml
  self-update: https://xxxxxx/ZetoHkr/update-config/raw/branch/main/update.yml

```

6. 在该程序的所在目录中再添加一个 update-helper.bat 更新辅助脚本即可完成自动检测更新并升级
7. 最终本地目录结构如下:
```text
.
├── config.yml
├── resource-update.exe
├── update-helper.bat
```

## 鸣谢以及用到的库

+ github.com/go-git/go-git/v5 用于实现Git的基本操作
+ github.com/mattn/go-colorable 解决Windows下CMD日志颜色乱码
+ gopkg.in/yaml.v2 解析YAML配置文件

> GoLand是一个IDE，可以最大限度地提高Go语言开发人员的工作效率。

特别感谢JetBrains为开源和教育学习提供免费的GoLand许可证。

[<img src="./jetbrains-variant-3.png" width="200"/>]()
