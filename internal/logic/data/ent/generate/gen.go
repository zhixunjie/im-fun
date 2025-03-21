package generate

import (
	"github.com/zhixunjie/im-fun/pkg/gomysql"
	"gorm.io/gen"
)

// DB 调试：tail -f /var/log/mysql/general_query.log
var DB, _ = gomysql.CreatePool(&gomysql.Config{
	Addr:     "127.0.0.1",
	Port:     "3306",
	UserName: "root",
	Password: "",
	Database: "im",
})

// Querier Dynamic SQL
// https://gorm.io/zh_CN/gen/dynamic_sql.html
// 定义一些通用的查询SQL：每个数据表在生成代码时，都会生成以下SQL的查询方法
type Querier interface {
	// GetByID SELECT * FROM @@table WHERE id=@id
	GetByID(id int64) (*gen.T, error) // GetByID query data by id and return it as *struct*
}

// 封装函数：根据表名生成代码
// 依赖函数：GenerateModel：生成对应的model文件和query文件
// https://gorm.io/zh_CN/gen/database_to_structs.html
func applyTableNames(g *gen.Generator) {
	g.ApplyInterface(
		func(Querier) {},
		g.GenerateModel("contact"),
		g.GenerateModel("message"),
		g.GenerateModel("user"),
		g.GenerateModel("robot"),
		g.GenerateModel("chat_group"),
		g.GenerateModel("chat_group_user"),
	)

}

// 封装函数使用：生成query文件和model文件
func genCode() {
	// gen.Config: https://gorm.io/zh_CN/gen/dao.html#gen-Config
	g := gen.NewGenerator(gen.Config{
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		OutPath:       "./query",
		ModelPkgPath:  "./model",
		FieldSignable: true,
	})

	g.UseDB(DB) // reuse your gorm db

	//根据表名生成代码
	applyTableNames(g)

	// Generate the code
	g.Execute()
}
