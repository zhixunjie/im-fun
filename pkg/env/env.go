package env

var (
	env     EnvType
	svcName string
)

type EnvType string

const (
	EnvTypeLocal EnvType = "local"
	EnvTypeTest  EnvType = "test"
	EnvTypeProd  EnvType = "prod"
)

// InitEnv 注入变量值
func InitEnv(e EnvType, svc string) {
	env = e
	svcName = svc
}

func IsLocal() bool {
	return env == EnvTypeLocal
}

func IsTest() bool {
	return env == EnvTypeTest
}

func IsProd() bool {
	return env == EnvTypeProd
}
