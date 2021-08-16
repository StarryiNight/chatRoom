package logic

import (
	"chatRoom/dao/mysql"
	"chatRoom/models"
	"chatRoom/pkg/jwt"
	"chatRoom/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}

	//雪花算法生成id
	userID := snowflake.GenID()

	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//保存数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return "", err
	}
	//生成自己定义的JWT
	return jwt.GenToken(user.UserID, user.Username)
}
