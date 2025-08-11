package controllers

import (
	"net/http"

	"github.com/qq272446181/BlogWebApi/config"
	"github.com/qq272446181/BlogWebApi/models"

	"github.com/gin-gonic/gin"
)

// 输入文字模型
type PostInput struct {
	//文章标题
	Title string `json:"title" binding:"required"`
	//文章内容
	Content string `json:"content" binding:"required"`
}

// 创建文章
// @Summary 创建文章
// @Description 创建文章
// @Tags 文章
// @Accept  json
// @Produce  json
// @Param input body PostInput true "文章输入模型"
// @Success 200 {object} models.ApiResponse
// @Router /posts [post]
func CreatePost(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	user, _ := c.Get("user")
	author := user.(models.User)

	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		res.Message = "输入格式错误"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	post := models.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: author.ID,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		res.Message = "创建文章失败"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Data = post
	res.Status = true
	res.Message = "创建文章成功"
	c.JSON(http.StatusCreated, res)
}

// 获取文章列表
// @Summary 获取文章列表
// @Description 获取文章列表
// @Tags 文章
// @Produce  json
// @Success 200 {object} models.ApiResponse
// @Router /posts [get]
func GetPosts(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	var posts []models.Post
	if err := config.DB.Preload("Author").Preload("Comments").Find(&posts).Error; err != nil {
		res.Message = "获取文章列表失败"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Data = posts
	res.Status = true
	res.Message = "获取文章列表成功"
	c.JSON(http.StatusOK, res)
}

// @Summary 根据ID获取文章
// @Description 根据ID获取文章
// @Tags 文章
// @Produce  json
// @Param id path int true "文章ID"
// @Success 200 {object} models.ApiResponse
// @Router /posts/{id} [get]
func GetPost(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	var post models.Post
	err := config.DB.Preload("Author").
		Preload("Comments").
		Preload("Comments.User").
		First(&post, c.Param("id")).Error
	if err != nil {
		res.Message = "获取文章失败"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Data = post
	res.Status = true
	res.Message = "获取文章成功"
	c.JSON(http.StatusOK, res)
}

// 更新文章
// @Summary 更新文章
// @Description 更新文章
// @Tags 文章
// @Accept  json
// @Produce  json
// @Param id path int true "文章ID"
// @Param input body PostInput true "文章输入模型"
// @Success 200 {object} models.ApiResponse
// @Router /posts/{id} [put]
func UpdatePost(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	user, _ := c.Get("user")
	currentUser := user.(models.User)

	var post models.Post
	if err := config.DB.First(&post, c.Param("id")).Error; err != nil {
		res.Message = "文章未找到"
		c.JSON(http.StatusNotFound, res)
		return
	}

	if post.AuthorID != currentUser.ID {
		res.Message = "您无权修改不属于自己的文章"
		c.JSON(http.StatusForbidden, res)
		return
	}

	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		res.Message = "输入格式错误"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	post.Title = input.Title
	post.Content = input.Content

	if err := config.DB.Save(&post).Error; err != nil {
		res.Message = "更新文章失败"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Data = post
	res.Status = true
	res.Message = "更新文章成功"
	c.JSON(http.StatusOK, res)
}

// 删除文章
// @Summary 删除文章
// @Description 删除文章
// @Tags 文章
// @Produce  json
// @Param id path int true "文章ID"
// @Success 200 {object} models.ApiResponse
// @Router /posts/{id} [delete]
func DeletePost(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	user, _ := c.Get("user")
	currentUser := user.(models.User)

	var post models.Post
	if err := config.DB.First(&post, c.Param("id")).Error; err != nil {
		res.Message = "文章未找到"
		c.JSON(http.StatusNotFound, res)
		return
	}

	if post.AuthorID != currentUser.ID {
		c.JSON(http.StatusForbidden, models.ApiResponse{
			Message: "您无权删除不属于自己的文章",
			Status:  false,
		})
		return
	}

	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusForbidden, models.ApiResponse{
			Message: "删除文章失败:" + err.Error(),
			Status:  false,
		})
		return
	}

	c.JSON(http.StatusForbidden, models.ApiResponse{
		Message: "删除文章成功",
		Status:  true,
	})
}
