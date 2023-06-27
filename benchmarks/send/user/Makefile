build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o send_user main.go
build-mac-m1:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o send_user main.go
build-mac-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o send_user main.go

run:
	./send_user --addr=localhost:8080 --start=1 --num=100000 --duration=30
	#./send_user --addr=localhost:8080 --start=1 --num=1000 --duration=30