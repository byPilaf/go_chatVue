package utils

import (
	"webchat/models"

	_ "github.com/go-sql-driver/mysql" //mysql驱动
	"github.com/jinzhu/gorm"
)

//DbMysql mysql数据库
var DbMysql *gorm.DB

func init() {
	//链接数据库
	DbMysql, _ = gorm.Open("mysql", "root:root@/webchat?charset=utf8&parseTime=True&loc=Local")

	//数据迁移
	DbMysql.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.User{})
}

//CloseMysqlDb 关闭数据库
func CloseMysqlDb() {
	DbMysql.Close()
}
