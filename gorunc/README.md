## CRI接口开发
模拟crictl的功能，实现对容器的管理，包括容器的创建、删除、查询、启动、停止、重启、进入容器、容器日志查看等功能。

一些基本功能
```bash
# 直接把代码拉下来，源码运行

go run main.go version   # 查看版本

go run main.go runp ./test/sandbox.yaml # 创建sandbox沙箱

go run main.go images  # 查看镜像列表

go run main.go run PodSandboxId ./test/ngx.yaml ./test/sandbox.yaml # 创建 单容器POD

go run main.go ps  # 查看容器列表

go run main.go exec -t [ContainerId] /bin/bash # 进入容器

```