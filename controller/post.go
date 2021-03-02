package controller

import (
	"sprout_server/common/constants"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/post"
	"sprout_server/models"
	"sprout_server/models/queryfields"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

func (pc *PostController) Create(c *gin.Context) {
	// 1. verify params
	var p models.ParamsAddPost
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	// 2. logic handle
	uid, exist := c.Get(constants.CtxUidKey)
	if !exist {
		response.Send(c, code.CodeNeedLogin)
	}
	statusCode := post.Create(&p, uid.(string))

	// 3. response result
	response.Send(c, statusCode)
	return

}

func (pc *PostController) Update(c *gin.Context) {
	var p models.ParamsUpdatePost
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	var u models.UriUpdatePost
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := post.Update(&p, &u)

	response.Send(c, statusCode)
	return
}

func (pc *PostController) GetPostList(c *gin.Context) {
	var qs models.QueryStringGetPostList
	if err := c.ShouldBindQuery(&qs); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	posts, statusCode := post.GetList(&qs)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, posts)
}

func (pc *PostController) GetTopPost(c *gin.Context) {
	topPost, statusCode := post.GetTopPost()

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, topPost)
}

func (pc *PostController) GetPostListByAdmin(c *gin.Context) {
	var queryFields queryfields.PostQueryFields
	if err := c.ShouldBindQuery(&queryFields); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	posts, statusCode := post.GetListByAdmin(&queryFields)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, posts)
}

func (pc *PostController) GetPostDetailByAdmin(c *gin.Context) {
	var p models.UriGetPostDetail
	if err := c.ShouldBindUri(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	postDetail, statusCode := post.GetDetailByAdmin(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, postDetail)
}

func (pc *PostController) GetPostDetail(c *gin.Context) {
	var p models.UriGetPostDetail
	if err := c.ShouldBindUri(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	postDetail, statusCode := post.GetDetail(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}

	response.SendWithData(c, statusCode, postDetail)
}
