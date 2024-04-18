package registry

import (
	"context"
	kratos_etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kratos_registry "github.com/go-kratos/kratos/v2/registry"
	"github.com/google/uuid"
	"github.com/zhixunjie/im-fun/pkg/host"
	"github.com/zhixunjie/im-fun/pkg/logging"
	etcdclient "go.etcd.io/etcd/client/v3"
	"net"
)

// kratos registry
// etcdctl-d get --prefix /microservices/

const EtcdAddr = "127.0.0.1:12379"

var KratosEtcdRegistry *kratos_etcd.Registry

func init() {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{EtcdAddr},
	})
	if err != nil {
		panic(err)
		return
	}
	// 基于etcd client，生成kratos的注册中心
	KratosEtcdRegistry = kratos_etcd.New(client)
}

type KratosServiceInstance struct {
	ServiceInstance    *kratos_registry.ServiceInstance
	GrpcServerListener net.Listener
}

// BuildServiceInstance 构建一个实例
func BuildServiceInstance(srvName, network, addr string) (*KratosServiceInstance, error) {
	// 创建监听请
	listener, err := net.Listen(network, addr)
	if err != nil {
		panic(err)
	}

	// 通过监听器，提取IP地址
	newAddr, err := host.Extract(addr, listener)
	if err != nil {
		logging.Errorf("Extract error: %v", err)
		return nil, err
	}

	instance := &KratosServiceInstance{
		ServiceInstance: &kratos_registry.ServiceInstance{
			ID:        uuid.NewString(),
			Name:      srvName,
			Endpoints: []string{"grpc://" + newAddr},
		},
		GrpcServerListener: listener,
	}

	return instance, nil
}

// Register 注册一个实例，并且返回其注销函数
func Register(ctx context.Context, instance *KratosServiceInstance) (func(), error) {
	// register to etcd
	err := KratosEtcdRegistry.Register(ctx, instance.ServiceInstance)
	if err != nil {
		logging.Errorf("Register error: %v", err)
		return nil, err
	}

	// 注销
	fn := func() {
		err = KratosEtcdRegistry.Deregister(ctx, instance.ServiceInstance)
		if err != nil {
			logging.Errorf("Deregister error: %v", err)
			return
		}
	}

	return fn, nil
}
