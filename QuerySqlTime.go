package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User 用户信息
type User struct {
	ID   int    `gorm:"primary_key;AUTO_INCREMENT;column:Id;"`
	Name string `gorm:"type:varchar(255);column:Name"`
}

// TableName 表名
func (User) TableName() string {
	return "t_user"
}

// 执行查询操作前
func queryBefore(scope *gorm.Scope) {
	start := time.Now()
	fmt.Println("start exec:", start)
	scope.Set("start", start)
}

// 执行查询操作后
func queryAfter(scope *gorm.Scope) {
	val, _ := scope.Get("start")
	start := val.(time.Time)
	end := time.Now()
	fmt.Println("end exec:", end)

	fmt.Println("执行耗时:", end.Sub(start))
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db.LogMode(true)
	db.AutoMigrate(&User{})

	// 注册
	db.Callback().Query().Before("gorm:query").Register("queryBefore", queryBefore)
	db.Callback().Query().After("gorm:query").Register("queryAfter", queryAfter)
}

func main() {
	defer db.Close()

	// db.Save(&User{Name: "test"})
	var user User
	db.First(&user)
	fmt.Println(user)
}
