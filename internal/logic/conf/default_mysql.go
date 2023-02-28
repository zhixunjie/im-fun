package conf

func DefaultMySQL() []MySQL {
	return []MySQL{
		{
			Addr:     "127.0.0.1:3306",
			UserName: "root",
			Password: "",
			DbName:   "im",
		},
	}
}

type MySQL struct {
	Addr     string
	UserName string
	Password string
	DbName   string
}
