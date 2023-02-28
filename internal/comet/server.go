package comet

import (
	"github.com/zhenjl/cityhash"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"math/rand"
	"time"
)

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
)

// Server 服务器（主体入口）
type Server struct {
	conf        *conf.Config // config
	round       *Round       // round sth
	buckets     []*Bucket    // bucket 数组
	bucketTotal uint32       // bucket总数

	//rpcClient logic.LogicClient
}

// NewServer returns a new Server.
func NewServer(conf *conf.Config) *Server {
	s := &Server{
		conf:  conf,
		round: NewRound(conf),
		//rpcClient: newLogicClient(conf.RPCClient),
	}

	// init bucket
	s.buckets = make([]*Bucket, conf.Bucket.Size)
	s.bucketTotal = uint32(conf.Bucket.Size)
	for i := 0; i < conf.Bucket.Size; i++ {
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

//func newLogicClient(conf *conf.RPCClient) logic.LogicClient {
//	conn, err := grpc.Dial("127.0.0.1:3119",
//		[]grpc.DialOption{
//			grpc.WithInsecure(),
//			grpc.WithInitialWindowSize(grpcInitialWindowSize),
//			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
//			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
//			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
//			grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
//			grpc.WithKeepaliveParams(keepalive.ClientParameters{
//				Time:                grpcKeepAliveTime,
//				Timeout:             grpcKeepAliveTimeout,
//				PermitWithoutStream: true,
//			}),
//		}...)
//	if err != nil {
//		panic(err)
//	}
//	return logic.NewLogicClient(conn)
//}
