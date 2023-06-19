# ARCH
ARCH=GOOS=linux GOARCH=amd64     # build target to intel64
#ARCH=GOOS=darwin GOARCH=arm64   # build target to m1

# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(ARCH) $(GOCMD) build
GOTEST=$(ARCH) $(GOCMD) test

#all: test build
all: build
build:
	rm -rf target/
	mkdir target/
	cp cmd/comet/comet.yaml target/comet.yaml
	cp cmd/logic/logic.yaml target/logic.yaml
	cp cmd/job/job.yaml target/job.yaml
	$(GOBUILD) -o target/comet cmd/comet/main.go
	$(GOBUILD) -o target/logic cmd/logic/main.go
	$(GOBUILD) -o target/job cmd/job/main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -rf target/

run:
	nohup target/logic -conf=target/logic.yaml -region=sh -zone=sh001 -deploy.env=dev -weight=10 2>&1 > target/logic.log &
	nohup target/comet -conf=target/comet.yaml -region=sh -zone=sh001 -deploy.env=dev -weight=10 -addrs=127.0.0.1 -debug=true 2>&1 > target/comet.log &
	nohup target/job -conf=target/job.yaml -region=sh -zone=sh001 -deploy.env=dev 2>&1 > target/job.log &

stop:
	pkill -f target/logic
	pkill -f target/job
	pkill -f target/comet
