# protocol
protoc --go_out=./proto --go_opt=paths=source_relative \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative protocol.proto

# logic
protoc --go_out=./pb --go_opt=paths=source_relative \
    --go-grpc_out=./pb  --go-grpc_opt=paths=source_relative logic.proto

# comet
protoc --go_out=./pb --go_opt=paths=source_relative \
    --go-grpc_out=./pb --go-grpc_opt=paths=source_relative comet.proto