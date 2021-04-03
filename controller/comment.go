package controller

import (
	"sprout_server/common/constants"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/comment"
	"sprout_server/models"
	"sprout_server/models/queryfields"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

func (cc *CommentController) CreatePostComment(c *gin.Context) {
	// 1. verify params
	var p models.ParamsAddComment
	if err := c.ShouldBindJSON(&p); err != nil {
		// params error
		// we use the shouldBindJSON and use the binding tag on model
		// the gin can help us to verify params
		response.Send(c, code.CodeInvalidParams)
		return
	}

	// get pid
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		response.Send(c, code.CodePostNotExist)
		return
	}
	p.Pid = pid
	// 2. logic handle
	// get uid from context
	uid, exist := c.Get(constants.CtxUidKey)
	if !exist {
		response.Send(c, code.CodeUserNotExist)
		return
	}
	p.Uid = uid.(string)

	ip, _ := c.Get(constants.CtxOriginIpKey)
	os, _ := c.Get(constants.CtxOriginOsKey)
	engine, _ := c.Get(constants.CtxOriginEngineKey)
	browser, _ := c.Get(constants.CtxOriginBrowserKey)
	if ip == nil {
		uid = ""
	}
	if os == nil {
		os = ""
	}
	if engine == nil {
		engine = ""
	}
	if browser == nil {
		browser = ""
	}
	statusCode := comment.CreatePostComment(&p, ip.(string), os.(string), engine.(string), browser.(string))

	// 3. response result
	response.Send(c, statusCode)
	return

}

func (cc *CommentController) AdminUpdatePostComment(c *gin.Context) {
	var p models.ParamsAdminUpdateComment
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	var u models.UriUpdateComment
	if err := c.ShouldBindUri(&u); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	statusCode := comment.AdminUpdatePostComment(&p, &u)

	response.Send(c, statusCode)
	return

}

func (cc *CommentController) GetPostCommentList(c *gin.Context) {
	var p models.ParamsGetCommentList
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	p.Pid = pid
	commentList, statusCode := comment.GetPostCommentList(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, commentList)
}

func (cc *CommentController) GetPostComments(c *gin.Context) {
	var p queryfields.CommentQueryFields
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	comments, statusCode := comment.GetPostComments(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, comments)
}

func (cc *CommentController) GetPublicComments(c *gin.Context) {
	var p queryfields.CommentQueryFields
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}

	p.ReviewStatus = 1
	p.IsDelete = 0

	comments, statusCode := comment.GetPostComments(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, comments)
}

func (cc *CommentController) GetPostParentCommentChildren(c *gin.Context) {
	var p models.ParamsGetParentCommentChildren
	if err := c.ShouldBindQuery(&p); err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	p.Pid = pid

	cid, err := strconv.ParseUint(c.Param("cid"), 10, 64)
	if err != nil {
		response.Send(c, code.CodeInvalidParams)
		return
	}
	p.Cid = cid
	parentComment, statusCode := comment.GetPostParentCommentChildren(&p)

	if statusCode != code.CodeOK {
		response.Send(c, statusCode)
		return
	}
	response.SendWithData(c, statusCode, parentComment)
}
