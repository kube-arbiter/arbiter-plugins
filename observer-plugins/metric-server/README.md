# metric-server
metric-server

## 注意问题
目前 `metric-server` 依赖 `arbiter`，但是 `arbieter.k8s.com.cn` 目前还不是代码仓库。
所以目前直接 `go mod` 是不可以的。需要这里的一个解决方法是利用 `go workspace`

```shell
mkdir arbiter.k8s.com.cn;
cd arbiter.k8s.com.cn;
git clone https://gitlab.dev.21vianet.com/arbiter/arbiter.git 
git clone https://gitlab.dev.21vianet.com/arbiter/arbiter-plugins.git

go work init
go work use ./arbiter ./arbiter-plugins/observer-plugins/metric-server
```

但是利用 `workspace` `go mod tidy` 命令还是无法正常工作[go mod tidy](https://github.com/golang/go/issues/50750)
如果要正常 `go mod tidy` 还需要  `replace`
整体目录如下

```shell
➜  arbiter.k8s.com.cn tree arbiter-plugins -L 4
arbiter-plugins
├── README.md
└── observer-plugins
    ├── README-CN.md
    ├── README.md
    ├── metric-server
    │   ├── Dockerfile
    │   ├── Makefile
    │   ├── README.md
    │   ├── example.yaml
    │   ├── go.mod
    │   ├── go.sum
    │   ├── main.go
    │   ├── metric_types.go
    │   └── server.go
    └── prometheus
        ├── Dockerfile
        ├── Makefile
        ├── go.mod
        ├── go.sum
        ├── main.go
        └── prometheus
            ├── prometheus.go
            └── server.go

4 directories, 19 files
➜  arbiter.k8s.com.cn tree . -L 1
.
├── arbiter
├── arbiter-plugins
├── go.work
└── go.work.sum

3 directories, 2 files
```