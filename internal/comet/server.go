package comet

import (
	"context"
	kratos_registry "github.com/go-kratos/kratos/v2/registry"
	"github.com/zhenjl/cityhash"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/micro_registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"math/rand"
	"time"
)

const (
	minServerHeartbeat = time.Minute * 41
	maxServerHeartbeat = time.Minute * 45
)

// Server 服务器（主体入口）
type Server struct {
	serviceInstance *kratos_registry.ServiceInstance
	conf            *conf.Config // config
	round           *Round       // round sth
	buckets         []*Bucket    // bucket 数组
	bucketTotal     uint32       // bucket总数

	rpcToLogic pb.LogicClient
}

// NewServer returns a new Server.
func NewServer(conf *conf.Config) *Server {
	s := &Server{
		serverId:   micro_registry.ServiceInstance,
		conf:       conf,
		round:      NewRound(conf),
		rpcToLogic: newLogicClient(conf.RPC.Client),
	}

	// init bucket
	s.buckets = make([]*Bucket, conf.Bucket.HashNum)
	s.bucketTotal = uint32(conf.Bucket.HashNum)
	for i := 0; i < conf.Bucket.HashNum; i++ {
		s.buckets[i] = NewBucket(conf.Bucket)
	}
	return s
}

// Buckets return all buckets.
func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

func (s *Server) AllocBucket(key string) *Bucket {
	idx := cityhash.CityHash32([]byte(key), uint32(len(key))) % s.bucketTotal

	return s.buckets[idx]
}

// RandHeartbeatTime 生成一个随机的心跳时间
func (s *Server) RandHeartbeatTime() time.Duration {
	return minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat)))
}

func (s *Server) Close() (err error) {
	return
}

// grpc options
const (
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 10
	grpcKeepAliveTimeout      = time.Second * 3
	grpcBackoffMaxDelay       = time.Second * 3
)

func newLogicClient(conf *conf.RPCClient) pb.LogicClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:12670",
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		}...)
	if err != nil {
		panic(err)
	}
	return pb.NewLogicClient(conn)
}
