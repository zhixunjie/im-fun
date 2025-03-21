package utils

import (
	"fmt"
	"gorm.io/gorm"
)

func CreateOrDrop(db *gorm.DB, op string, srcTbName string, shardingNum int64) {
	switch op {
	case "create":
		for i := int64(0); i < shardingNum; i++ {
			tbName := fmt.Sprintf("%s_%d", srcTbName, i)
			query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s LIKE %s", tbName, srcTbName)
			if err := db.Exec(query).Error; err != nil {
				fmt.Printf("❌ 创建表 %s 失败: %v\n", tbName, err)
				break
			}
		}
	case "drop":
		for i := int64(0); i < shardingNum; i++ {
			tbName := fmt.Sprintf("%s_%d", srcTbName, i)
			query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tbName)
			if err := db.Exec(query).Error; err != nil {
				fmt.Printf("❌ 删除表 %s 失败: %v\n", tbName, err)
				break
			}
		}
	}
}
