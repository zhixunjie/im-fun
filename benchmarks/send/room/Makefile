build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o send_room main.go
build-mac-m1:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o send_room main.go
build-mac-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o send_room main.go

run:
	./send_room --addr=localhost:8080 --start=1 --num=20000 --duration=10