# make 命令配置，被 ci.yml 引用，请勿删除
# 初始化项目目录变量
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output

GOPKGS  := $$(go list ./...| grep -vE "vendor")

# 设置编译时所需要的 Go 环境
export GOENV = $(HOMEDIR)/go.env

#执行编译，可使用命令 make 或 make all 执行， 顺序执行 prepare -> compile -> test -> package 几个阶段
all: prepare compile test package

# prepare阶段， 使用 bcloud 下载非 Go 依赖，可单独执行命令: make prepare
prepare: prepare-dep
prepare-dep:
	#bcloud local -U # 下载非 Go 依赖，依赖之前的 BCLOUD 文件
	git version     # 低于 2.17.1 可能不能正常工作
	go env          # 打印出 go 环境信息，可用于排查问题

set-env:
	go mod download -x || go mod download -x # 下载 Go 依赖

# compile 阶段，执行编译命令，可单独执行命令: make compile
compile:build
build: set-env
	echo "build: do nothing"
	mkdir -p $(OUTDIR)/build
	#go build -o $(OUTDIR)/build main.go

# test 阶段，进行单元测试， 可单独执行命令: make test
# cover 平台会优先执行此命令
test: test-case
test-case: set-env
	mkdir -p $(OUTDIR)/test
	go test -race -timeout=120s -v -cover $(GOPKGS) -coverprofile=$(OUTDIR)/test/coverage.out | tee $(OUTDIR)/test/unittest.txt

# package 阶段，对编译产出进行打包，输出到 output 目录， 可单独执行命令: make package
package: package-bin
package-bin:
	echo "package: do nothing"
	mkdir -p $(OUTDIR)/bin

# clean 阶段，清除过程中的输出， 可单独执行命令: make clean
clean:
	rm -rf $(OUTDIR)

# avoid filename conflict and speed up build
.PHONY: all prepare compile test package  clean build