package controllers

import (
	"net/http"

	"github.com/qq272446181/BlogWebApi/config"
	"github.com/qq272446181/BlogWebApi/models"

	"github.com/gin-gonic/gin"
)

// 评论输入模型
type CommentInput struct {
	//评论内容
	Content string `json:"content" binding:"required"`
}

// CreateComment 新增评论
// @Summary 新增评论
// @Description 新增评论
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param id path int true "文章ID"
// @Param comment body CommentInput true "评论内容"
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 404 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /posts/{id}/addcomment [post]
func CreateComment(c *gin.Context) {
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

	var input CommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		res.Message = "输入参数错误"
		c.JSON(http.StatusNotFound, res)
		return
	}

	comment := models.Comment{
		Content: input.Content,
		PostID:  post.ID,
		UserID:  currentUser.ID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		res.Message = "文章保存失败：" + err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Status = true
	res.Data = comment
	c.JSON(http.StatusCreated, res)
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 删除评论
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param id path int true "文章ID"
// @Param commentId path int true "评论ID"
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 404 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /posts/{id}/deletecomment/{commentId} [delete]
func DeleteComment(c *gin.Context) {
	res := models.ApiResponse{
		Status: false,
	}
	user, _ := c.Get("user")
	currentUser := user.(models.User)

	var comment models.Comment
	if err := config.DB.First(&comment, c.Param("commentId")).Error; err != nil {
		res.Message = "评论未找到"
		c.JSON(http.StatusNotFound, res)
		return
	}

	if comment.UserID != currentUser.ID {
		res.Message = "您没有权限删除该评论"
		c.JSON(http.StatusForbidden, res)
		return
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		res.Message = "评论删除失败：" + err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Status = true
	res.Message = "评论删除成功"
	c.JSON(http.StatusOK, res)
}
