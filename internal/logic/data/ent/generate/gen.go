package generate

import (
	"fmt"
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
		g.GenerateModel("group"),
		g.GenerateModel("group_user"),
	)

}

// 封装函数使用：生成query文件和model文件
func genCode2() {
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

func genCreateTable() {
	tableNum := 100
	for i := 0; i < tableNum; i++ {
		fmt.Printf("create table `message_%v` like `message`;\n", i)
	}
	for i := 0; i < tableNum; i++ {
		fmt.Printf("create table `contact_%v` like `contact`;\n", i)
	}
}

/**
create table `message_0` like `message`;
create table `message_1` like `message`;
create table `message_2` like `message`;
create table `message_3` like `message`;
create table `message_4` like `message`;
create table `message_5` like `message`;
create table `message_6` like `message`;
create table `message_7` like `message`;
create table `message_8` like `message`;
create table `message_9` like `message`;
create table `message_10` like `message`;
create table `message_11` like `message`;
create table `message_12` like `message`;
create table `message_13` like `message`;
create table `message_14` like `message`;
create table `message_15` like `message`;
create table `message_16` like `message`;
create table `message_17` like `message`;
create table `message_18` like `message`;
create table `message_19` like `message`;
create table `message_20` like `message`;
create table `message_21` like `message`;
create table `message_22` like `message`;
create table `message_23` like `message`;
create table `message_24` like `message`;
create table `message_25` like `message`;
create table `message_26` like `message`;
create table `message_27` like `message`;
create table `message_28` like `message`;
create table `message_29` like `message`;
create table `message_30` like `message`;
create table `message_31` like `message`;
create table `message_32` like `message`;
create table `message_33` like `message`;
create table `message_34` like `message`;
create table `message_35` like `message`;
create table `message_36` like `message`;
create table `message_37` like `message`;
create table `message_38` like `message`;
create table `message_39` like `message`;
create table `message_40` like `message`;
create table `message_41` like `message`;
create table `message_42` like `message`;
create table `message_43` like `message`;
create table `message_44` like `message`;
create table `message_45` like `message`;
create table `message_46` like `message`;
create table `message_47` like `message`;
create table `message_48` like `message`;
create table `message_49` like `message`;
create table `message_50` like `message`;
create table `message_51` like `message`;
create table `message_52` like `message`;
create table `message_53` like `message`;
create table `message_54` like `message`;
create table `message_55` like `message`;
create table `message_56` like `message`;
create table `message_57` like `message`;
create table `message_58` like `message`;
create table `message_59` like `message`;
create table `message_60` like `message`;
create table `message_61` like `message`;
create table `message_62` like `message`;
create table `message_63` like `message`;
create table `message_64` like `message`;
create table `message_65` like `message`;
create table `message_66` like `message`;
create table `message_67` like `message`;
create table `message_68` like `message`;
create table `message_69` like `message`;
create table `message_70` like `message`;
create table `message_71` like `message`;
create table `message_72` like `message`;
create table `message_73` like `message`;
create table `message_74` like `message`;
create table `message_75` like `message`;
create table `message_76` like `message`;
create table `message_77` like `message`;
create table `message_78` like `message`;
create table `message_79` like `message`;
create table `message_80` like `message`;
create table `message_81` like `message`;
create table `message_82` like `message`;
create table `message_83` like `message`;
create table `message_84` like `message`;
create table `message_85` like `message`;
create table `message_86` like `message`;
create table `message_87` like `message`;
create table `message_88` like `message`;
create table `message_89` like `message`;
create table `message_90` like `message`;
create table `message_91` like `message`;
create table `message_92` like `message`;
create table `message_93` like `message`;
create table `message_94` like `message`;
create table `message_95` like `message`;
create table `message_96` like `message`;
create table `message_97` like `message`;
create table `message_98` like `message`;
create table `message_99` like `message`;
create table `contact_0` like `contact`;
create table `contact_1` like `contact`;
create table `contact_2` like `contact`;
create table `contact_3` like `contact`;
create table `contact_4` like `contact`;
create table `contact_5` like `contact`;
create table `contact_6` like `contact`;
create table `contact_7` like `contact`;
create table `contact_8` like `contact`;
create table `contact_9` like `contact`;
create table `contact_10` like `contact`;
create table `contact_11` like `contact`;
create table `contact_12` like `contact`;
create table `contact_13` like `contact`;
create table `contact_14` like `contact`;
create table `contact_15` like `contact`;
create table `contact_16` like `contact`;
create table `contact_17` like `contact`;
create table `contact_18` like `contact`;
create table `contact_19` like `contact`;
create table `contact_20` like `contact`;
create table `contact_21` like `contact`;
create table `contact_22` like `contact`;
create table `contact_23` like `contact`;
create table `contact_24` like `contact`;
create table `contact_25` like `contact`;
create table `contact_26` like `contact`;
create table `contact_27` like `contact`;
create table `contact_28` like `contact`;
create table `contact_29` like `contact`;
create table `contact_30` like `contact`;
create table `contact_31` like `contact`;
create table `contact_32` like `contact`;
create table `contact_33` like `contact`;
create table `contact_34` like `contact`;
create table `contact_35` like `contact`;
create table `contact_36` like `contact`;
create table `contact_37` like `contact`;
create table `contact_38` like `contact`;
create table `contact_39` like `contact`;
create table `contact_40` like `contact`;
create table `contact_41` like `contact`;
create table `contact_42` like `contact`;
create table `contact_43` like `contact`;
create table `contact_44` like `contact`;
create table `contact_45` like `contact`;
create table `contact_46` like `contact`;
create table `contact_47` like `contact`;
create table `contact_48` like `contact`;
create table `contact_49` like `contact`;
create table `contact_50` like `contact`;
create table `contact_51` like `contact`;
create table `contact_52` like `contact`;
create table `contact_53` like `contact`;
create table `contact_54` like `contact`;
create table `contact_55` like `contact`;
create table `contact_56` like `contact`;
create table `contact_57` like `contact`;
create table `contact_58` like `contact`;
create table `contact_59` like `contact`;
create table `contact_60` like `contact`;
create table `contact_61` like `contact`;
create table `contact_62` like `contact`;
create table `contact_63` like `contact`;
create table `contact_64` like `contact`;
create table `contact_65` like `contact`;
create table `contact_66` like `contact`;
create table `contact_67` like `contact`;
create table `contact_68` like `contact`;
create table `contact_69` like `contact`;
create table `contact_70` like `contact`;
create table `contact_71` like `contact`;
create table `contact_72` like `contact`;
create table `contact_73` like `contact`;
create table `contact_74` like `contact`;
create table `contact_75` like `contact`;
create table `contact_76` like `contact`;
create table `contact_77` like `contact`;
create table `contact_78` like `contact`;
create table `contact_79` like `contact`;
create table `contact_80` like `contact`;
create table `contact_81` like `contact`;
create table `contact_82` like `contact`;
create table `contact_83` like `contact`;
create table `contact_84` like `contact`;
create table `contact_85` like `contact`;
create table `contact_86` like `contact`;
create table `contact_87` like `contact`;
create table `contact_88` like `contact`;
create table `contact_89` like `contact`;
create table `contact_90` like `contact`;
create table `contact_91` like `contact`;
create table `contact_92` like `contact`;
create table `contact_93` like `contact`;
create table `contact_94` like `contact`;
create table `contact_95` like `contact`;
create table `contact_96` like `contact`;
create table `contact_97` like `contact`;
create table `contact_98` like `contact`;
create table `contact_99` like `contact`;
*/
