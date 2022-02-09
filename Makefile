VERSION = 1.0.0
BUILDTIME=`date +%FT%T%z`
REPOURL=`git config --get remote.origin.url | sed 's/:/\:/' | sed 's/git@/https:\/\//' | sed 's/\.git//'`
LDFLAGS = -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILDTIME} -X main.RepoURL=${REPOURL}"
GOARCH = amd64
APP := $(shell basename $(CURDIR))
BUILDPLATFORM = darwin

linux: export GOOS=linux
darwin: export GOOS=darwin
windows: export GOOS=windows

all: modules linux darwin windows

linux:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/${APP}-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/${APP}_${VERSION}-${GOOS}_${GOARCH}.zip ${APP}
	rm -f ${APP}

darwin:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/${APP}${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/${APP}_${VERSION}-${GOOS}_${GOARCH}.zip ${APP}
	rm -f ${APP}

windows:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/${APP}${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/${APP}_${VERSION}-${GOOS}_${GOARCH}.zip ${APP}.exe
	rm -f ${APP}.exe

modules:
	rm -f go.mod
	go mod init ${APP} 	
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/credentials
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/iam
	go get github.com/aws/aws-sdk-go/service/sts

.PHONY: clean
clean:
	rm -rf release
	rm -f ${APP} ${APP}.exe
	rm -f go.mod go.sum
