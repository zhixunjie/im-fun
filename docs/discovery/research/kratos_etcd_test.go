package research

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func TestKratosEtcd(t *testing.T) {
	srvName := "hello.world"
	version := "1.2.10"
	hs := http.NewServer()
	gs := grpc.NewServer()

	// new etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	// new reg with etcd client
	reg := etcd.New(client)

	app := kratos.New(
		// service-name
		kratos.Name(srvName),
		kratos.Version(version),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			hs,
			gs,
		),
		// with registrar
		kratos.Registrar(reg),
	)
	app.Run()
}
