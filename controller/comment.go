package controller

import (
	"sprout_server/common/constants"
	"sprout_server/common/response"
	"sprout_server/common/response/code"
	"sprout_server/logic/comment"
	"sprout_server/models"
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
	statusCode := comment.CreatePostComment(&p)

	// 3. response result
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
