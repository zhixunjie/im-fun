package comet

import (
	"github.com/zhenjl/cityhash"
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/registry"
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
