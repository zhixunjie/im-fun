build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tcp-client main.go
build-mac-m1:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o tcp-client main.go
build-mac-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o tcp-client main.go

run:
	./tcp-client --addr=localhost:12571 --start=1 --num=1