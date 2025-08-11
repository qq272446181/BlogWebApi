package controllers

import (
	"net/http"
	"time"

	"github.com/qq272446181/BlogWebApi/config"
	"github.com/qq272446181/BlogWebApi/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthInput 输入用户模型
// @Description 用户账户信息
type AuthInput struct {
	//用户账号
	Username string `json:"username" binding:"required" example:"admin"`
	//密码
	Password string `json:"password" binding:"required" example:"password123"`
	Email    string `json:"email" example:"272@qq.com"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Register godoc
// @Summary 用户注册
// @Description 注册新用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param input body AuthInput true "注册信息"
// @Success 201 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	var input AuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	user := models.User{Username: input.Username, Email: input.Email}

	err := user.HashPassword(input.Password)
	if err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var icount int64
	err = config.DB.Model(&models.User{}).Where("username = ?", input.Username).Count(&icount).Error
	if err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}
	if icount > 0 {
		res.Message = "用户名已存在"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}
	res.Status = true
	res.Message = "用户注册成功"
	c.JSON(http.StatusCreated, res)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录获取JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param input body AuthInput true "登录信息"
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	var input AuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {

		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {

		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	expirationTime := time.Now().Add(time.Duration(config.GetConfig().JWT.ExpireHours) * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}
	res.Status = true
	res.Message = "登录成功"
	res.Data = tokenString

	c.JSON(http.StatusOK, res)
}
