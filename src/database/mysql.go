package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"open-devops/src/modules/server/config"
	"os"
	"time"
)

var db *gorm.DB

func InitMysql(mysqlS *config.MySQLConf) {
	var err error
	dsn := mysqlS.DB

	newLogger := logger.New(
		log.New(os.Stdout, "rn", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * 3, // 慢 SQL 阈值
			LogLevel:      logger.Info,     //logger.Silent //不进行任何打印
			Colorful:      true,            // 色彩打印
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键
	})
	if err != nil {
		panic(fmt.Sprintf("MysqlConnectError: %s", err.Error()))
	} else {
		if sqlDb, err := db.DB(); err != nil {
			panic(fmt.Sprintf("MysqlConnectError: " + err.Error()))
		} else {
			sqlDb.SetMaxIdleConns(mysqlS.Idle) //保留的空间链接最大数
			sqlDb.SetMaxOpenConns(mysqlS.Max)  //最大开启的链接数
			// sqlDb.SetConnMaxIdleTime(time.Second * 15) //空闲连接持续一定时间后关闭
			sqlDb.SetConnMaxLifetime(time.Second * 15) //arm 环境下上条【连接池里最大空闲连接数】编辑报错 ， 使用此行【连接池里面的连接最大空闲时长】
		}

	}
	fmt.Println("MysqlConnectSuccess")
}

// GetDb - get a database connection
func GetDb() *gorm.DB {
	//log.Info(fmt.Sprintf("%#v", db))
	return db
}
