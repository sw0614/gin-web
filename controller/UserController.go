package controller

import (
	"ginpro/common"
	"ginpro/model"
	"ginpro/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": 500,
			"msg":  "加密错误",
		})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)
	//返回结果
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "用户注册成功",
	})

}
func Login(ctx *gin.Context) {

	DB := common.GetDB()

	//name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"code": 422,
			"msg":  "手机号必须为11位",
		})
		return
	}
	//判断用户是否存在
	var user model.User
	DB.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"code": 422,
			"msg":  "用户不存在",
		})
		return
	}

	if !isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"code": 422,
			"msg":  "用户已存在",
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

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	}

	token := "11"

	//返回结果
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{"token": token},
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
