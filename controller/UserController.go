package controller

import (
	"ginpro/common"
	"ginpro/model"
	"ginpro/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()

	//获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"code": 422,
			"msg":  "手机号必须为11位",
		})
		return
	}

	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"code": 422,
			"msg":  "密码必须大于6位",
		})
		return
	}
	//如果名称没有传，给一个10位的字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	//判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"code": 422,
			"msg":  "用户已存在",
		})
		return
	}
	//创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)
	//返回结果
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "用户注册成功",
	})

}
func Login(ctx *gin.Engine) {
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
