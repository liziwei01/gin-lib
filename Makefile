#初始化项目目录变量
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output

APPNAME := $(shell basename `pwd`)
OUTPUT_FILE := ${APPNAME}.tar.gz

# GOROOT  := /usr/local
# GOROOT  := /home/linuxbrew/.linuxbrew
# GO      := $(GOROOT)/bin/go
GO      := $(shell which go)
GOPATH  := $(shell $(GO) env GOPATH)
GOMOD   := $(GO) mod
GOBUILD := $(GO) build
GOTEST  := $(GO) test -gcflags="-N -l"
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")

# make, make all
all: prepare compile package

# make prepare, download dependencies
prepare: gomod
gomod: prepare-dep set-env 
	$(GOMOD) tidy

prepare-dep:
#	git config --global http.sslVerify false

set-env:
	$(GO) env -w GOPROXY="https://proxy.golang.org,direct"
#	$(GO) env -w GOPROXY="https://goproxy.io,direct"

# make compile
compile: build
build:
	$(GOBUILD) -o $(HOMEDIR)/bin/$(APPNAME)

# make test, test your code
test: prepare test-case
test-case:
	$(GOTEST) -v -cover $(GOPKGS)

package: package-bin
package-bin:
	$(shell rm -rf $(OUTDIR))
	$(shell mkdir -p $(OUTDIR))
	$(shell cp -a bin $(OUTDIR)/bin)
	$(shell cp -a conf $(OUTDIR)/conf)
	$(shell cd $(OUTDIR)/; tar -zcf ${OUTPUT_FILE} ./*; rm -rf bin conf)

# clean阶段, 清除过程中的输出, 可单独执行命令: make clean
clean:
	rm -rf $(OUTDIR) $(HOMEDIR)/bin