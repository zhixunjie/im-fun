package research

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestKratosConsul(t *testing.T) {
	srvName := "hello.world"
	version := "1.2.10"
	hs := http.NewServer()
	gs := grpc.NewServer()

	// new consul client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	// new reg with consul client
	reg := consul.New(client)

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
