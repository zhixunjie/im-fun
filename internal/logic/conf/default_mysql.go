package conf

func DefaultMySQL() []MySQL {
	return []MySQL{
		{
			Addr:     "127.0.0.1",
			Port:     "3306",
			UserName: "root",
			Password: "",
			DbName:   "im",
		},
	}
}

type MySQL struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbName"`
}
