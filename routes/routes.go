package routes

import (
	"sprout_server/common/myvali"
	"sprout_server/controller"
	"sprout_server/logger"
	"sprout_server/middlewares"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func Setup() (*gin.Engine, error) {
	r := gin.New()
	if err := myvali.Init(); err != nil {
		return r, err
	}
	r.Use(logger.GinLogger(), logger.GinRecovery(true), cors.New(cors.Config{
		AllowOrigins:     []string{"https://init.center", "http://init.center", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"content-type", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// register the routes

	// create the user route group
	user := r.Group("/users")
	{
		userController := &controller.UserController{}
		user.POST("", userController.SignUp)

	}

	session := r.Group("/session")
	{
		sessionController := &controller.SessionController{}
		session.POST("", sessionController.SignIn)
		session.GET("", middlewares.JwtAuth(), sessionController.CheckSignIn)
		session.GET("/admin", middlewares.JwtAuth(), middlewares.AdminAuth(), sessionController.CheckSignIn)
	}

	vCode := r.Group("/vcode")
	{
		vCodeController := &controller.VCodeController{}
		vCode.POST("/ecode", vCodeController.SendECode)
	}

	category := r.Group("/categories")
	{
		categoryController := &controller.CategoryController{}
		category.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), categoryController.Create)
		category.PUT("/:id", middlewares.JwtAuth(), middlewares.AdminAuth(), categoryController.Update)
		category.DELETE("/:id", middlewares.JwtAuth(), middlewares.AdminAuth(), categoryController.Delete)
		category.GET("", categoryController.GetByQuery)
	}

	tag := r.Group("/tags")
	{
		tagController := &controller.TagController{}
		tag.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), tagController.Create)
		tag.PUT("/:id", middlewares.JwtAuth(), middlewares.AdminAuth(), tagController.Update)
		tag.DELETE("/:id", middlewares.JwtAuth(), middlewares.AdminAuth(), tagController.Delete)
		tag.GET("", tagController.GetByQuery)
	}

	post := r.Group("/posts")
	{
		postController := &controller.PostController{}
		post.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), postController.Create)
		post.PUT("/:pid", middlewares.JwtAuth(), middlewares.AdminAuth(), postController.Update)
		post.GET("", postController.GetPostList)
		post.GET("/:pid", postController.GetPostDetail)
	}

	top := r.Group("/top")
	{
		postController := &controller.PostController{}
		top.GET("/post", postController.GetTopPost)
	}

	admin := r.Group("/admin", middlewares.JwtAuth(), middlewares.AdminAuth())
	{
		adminPost := admin.Group("/posts")
		{
			adminPostController := &controller.PostController{}
			adminPost.GET("", adminPostController.GetPostListByAdmin)
			adminPost.GET("/:pid", adminPostController.GetPostDetailByAdmin)

		}

		adminComment := admin.Group("/comments")
		{
			adminCommentController := &controller.CommentController{}
			adminComment.GET("/posts", adminCommentController.GetPostComments)
			adminComment.PUT("/:cid", adminCommentController.AdminUpdatePostComment)
		}

		adminUser := admin.Group("/users")
		{
			adminUserController := &controller.UserController{}
			adminUser.GET("", adminUserController.AdminGetUsers)
			adminUser.PUT("/:uid", adminUserController.AdminUpdateUser)
			adminUser.POST("/:uid/ban", adminUserController.BanUser)
			adminUser.DELETE("/:uid/ban", adminUserController.UnblockUser)
		}
	}

	comment := r.Group("/comments")
	{
		commentController := &controller.CommentController{}
		comment.GET("/posts/:pid", commentController.GetPostCommentList)
		// why not use /:cid/posts/:pid ?
		// because the gin(httpRouter) does not support it, it will cause conflicts and panic directly
		comment.GET("/posts/:pid/comment/:cid/children", commentController.GetPostParentCommentChildren)
		comment.POST("/posts/:pid", middlewares.JwtAuth(), commentController.CreatePostComment)
	}

	favorite := r.Group("/favorites")
	{
		favoriteController := &controller.FavoriteController{}
		favorite.GET("/posts/:pid", middlewares.JwtAuth(), favoriteController.CheckUserFavoritePost)
		favorite.POST("/posts/:pid", middlewares.JwtAuth(), favoriteController.AddUserFavoritePost)
		favorite.DELETE("/posts/:pid", middlewares.JwtAuth(), favoriteController.DeleteUserFavoritePost)
	}

	return r, nil
}
