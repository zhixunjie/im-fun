package env

var (
	env     Type
	svcName string
)

type Type string

const (
	TypeLocal Type = "local"
	TypeTest  Type = "test"
	TypeProd  Type = "prod"
)

// InitEnv 注入变量值
func InitEnv(e Type, svc string) {
	env = e
	svcName = svc
}

func IsLocal() bool {
	return env == TypeLocal
}

func IsTest() bool {
	return env == TypeTest
}

func IsProd() bool {
	return env == TypeProd
}
