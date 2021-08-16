package controllers

import (
	"chatRoom/dao/redis"
	"chatRoom/logic"
	"chatRoom/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的路由
func SignUpHandler(c *gin.Context) {
	//获取参数
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数错误",
			"err": err,
		})
		return
	}
	//校验是否为空
	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数错误",
		})
		return
	}
	fmt.Println(p)

	//业务处理
	if err := logic.SignUp(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
			"err": err,
		})
		return
	}

	//返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})

}

func LoginHandler(c *gin.Context) {
	//1.获取请求参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数错误",
			"err": err,
		})
		return
	}
	//2.业务逻辑处理
	token, err := logic.Login(p)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "登陆失败",
			"err": err,
		})
		return
	}
	//3.redis验证用户登陆是否唯一
	ttl, _ := redis.Rdb.Do("TTL", "key").Int64()
	fmt.Println(ttl)
	if ttl > 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "登陆失败",
			"err": "请勿重复登陆",
		})
		return
	}
	//未登陆 将用户名和token保存
	redis.Rdb.Set(p.Username, token, time.Hour*2)

	//4.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg":   "登陆成功",
		"token": token,
	})
}

func Logout(c *gin.Context) {
	username, exist := c.Get("username")
	if exist {
		redis.Rdb.Del(username.(string))
		c.JSON(http.StatusOK, gin.H{
			"msg": "注销成功",
		})
	}
}
