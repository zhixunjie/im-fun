package micro_registry

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kratos_registry "github.com/go-kratos/kratos/v2/registry"
	"github.com/google/uuid"
	"github.com/zhixunjie/im-fun/pkg/host"
	"github.com/zhixunjie/im-fun/pkg/logging"
	etcdclient "go.etcd.io/etcd/client/v3"
	"net"
)

func Register(ctx context.Context, srvName, addr string, listener net.Listener) (func(), error) {
	// get registry
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{"127.0.0.1:12379"},
	})
	if err != nil {
		logging.Errorf("BuildInstance error: %v", err)
		return nil, err
	}
	registry := etcd.New(client)

	// build instance
	instance, err := BuildInstance(srvName, addr, listener)
	if err != nil {
		logging.Errorf("BuildInstance error: %v", err)
		return nil, err
	}

	// register to etcd
	err = registry.Register(ctx, instance)
	if err != nil {
		logging.Errorf("Register error: %v", err)
		return nil, err
	}

	// 注销
	fn := func() {
		err = registry.Deregister(ctx, instance)
		if err != nil {
			logging.Errorf("Deregister error: %v", err)
			return
		}
	}

	return fn, nil
}

// BuildInstance 构建一个实例对象
func BuildInstance(name, address string, lis net.Listener) (*kratos_registry.ServiceInstance, error) {
	// 提取IP地址
	addr, err := host.Extract(address, lis)
	if err != nil {
		return nil, err
	}

	return &kratos_registry.ServiceInstance{
		ID:        uuid.NewString(),
		Name:      name,
		Endpoints: []string{"grpc://" + addr},
	}, nil
}
