package comet

import (
	"context"
	kratos_discovery "github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
	"github.com/zhenjl/cityhash"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"math/rand"
	"time"
)

const (
	minServerHeartbeat = time.Minute * 41
	maxServerHeartbeat = time.Minute * 45
)

// TcpServer 服务器（主体入口）
type TcpServer struct {
	serverId    string       // 服务器ID
	conf        *conf.Config // config
	round       *Round       // round sth
	buckets     []*Bucket    // bucket 数组
	bucketTotal uint32       // bucket总数

	rpcToLogic pb.LogicClient
}

// NewTcpServer returns a new TcpServer.
func NewTcpServer(conf *conf.Config) (*TcpServer, *registry.KratosServiceInstance) {
	// build grpc server instance
	rpcConfig := conf.RPC.Server
	instance, err := registry.BuildServiceInstance(conf.Name, rpcConfig.Network, rpcConfig.Addr)
	if err != nil {
		panic(err)
	}

	s := &TcpServer{
		serverId:   instance.ServiceInstance.ID,
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
	return s, instance
}

// Buckets return all buckets.
func (s *TcpServer) Buckets() []*Bucket {
	return s.buckets
}

func (s *TcpServer) AllocBucket(key string) *Bucket {
	idx := cityhash.CityHash32([]byte(key), uint32(len(key))) % s.bucketTotal

	return s.buckets[idx]
}

// RandHeartbeatTime 生成一个随机的心跳时间
func (s *TcpServer) RandHeartbeatTime() time.Duration {
	return minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat)))
}

func (s *TcpServer) Close() (err error) {
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

	conn, err := grpc.DialContext(ctx, "discovery:///logic",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
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
			// GRPC服务发现机制（使用Kratos的注册中心和解析器）
			grpc.WithResolvers(kratos_discovery.NewBuilder(registry.KratosEtcdRegistry, kratos_discovery.WithInsecure(true))),
		}...)
	if err != nil {
		panic(err)
	}
	return pb.NewLogicClient(conn)
}
