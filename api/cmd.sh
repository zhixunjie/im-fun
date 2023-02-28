protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative protocol/protocol.proto
protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative logic/logic.proto
protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative comet/comet.proto