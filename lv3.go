package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID       uint
	Name     string `gorm:"unique"`
	Password string
	Question string
	Answer   string
}
type Talk struct {
	ID       uint `gorm:"primary_key`
	Sender   string
	Reciver  string
	Dialogue string
}

func main() {
	r := gin.Default()
	db, err := gorm.Open("mysql", "root:root@/db1?charset=utf8mb4")
	if err != nil {
		print("gorm err:%v", err)
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Talk{})
	r.POST("/user", func(c *gin.Context) {
		name := c.PostForm("name")
		password := c.PostForm("password")
		question := c.PostForm("question")
		answer := c.PostForm("answer")
		u := User{
			Name:     name,
			Password: password,
			Question: question,
			Answer:   answer,
		}
		err := db.Create(&u)
		if err.Error != nil {
			fmt.Printf("register wrong: %v\n", err.Error)
		} else {
			c.String(200, "register!")
		}
	}) //注册账户
	r.GET("/user", func(c *gin.Context) {
		var user = new(User)
		name := c.Query("name")
		password := c.Query("password")
		op := c.Query("option")
		db.Where("name=?", name).Find(&user)
		if user.Password == password {
			c.String(200, "login!")
			if op == "1" {
				reciver := c.Query("reciver")
				message := c.Query("message")
				dialogue := Talk{
					Sender:   name,
					Reciver:  reciver,
					Dialogue: message,
				}
				err := db.Create(&dialogue)
				if err.Error != nil {
					c.String(200, "send failed", err.Error)
				} else {
					c.String(200, "send!")
				}
			}
		} else {
			c.String(200, "name or password wrong!")
		}
	}) //登录账户,op为1时发送留言
	r.PUT("/user", func(c *gin.Context) {
		name := c.PostForm("name")
		answer := c.PostForm("answer")
		newPassword := c.PostForm("newPassword")
		op := c.PostForm("option")
		var res User
		db.Where("name=?", name).First(&res)
		if res.Answer == answer {
			if op == "1" {
				c.String(200, name+"'password is:"+res.Password)
			} else if op == "2" {
				res.Password = newPassword
				db.Save(&res)
				c.String(200, "OK")
			}
		} else {
			c.String(200, "wrong answer")
		}
	}) //op为1找回密码。op为2修改密码
	defer db.Close()
	r.Run(":8080")
}

