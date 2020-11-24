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
		category.GET("", categoryController.GetAll)
	}

	tag := r.Group("/tags")
	{
		tagController := &controller.TagController{}
		tag.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), tagController.Create)
		tag.GET("", tagController.GetAll)
	}

	post := r.Group("/posts")
	{
		postController := &controller.PostController{}
		post.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), postController.Create)
		post.GET("", postController.GetPostList)
		post.GET("/:pid", postController.GetPostDetail)
	}

	comment := r.Group("/comments")
	{
		commentController := &controller.CommentController{}
		comment.GET("posts/:pid", commentController.GetPostCommentList)
		// why not use /:cid/posts/:pid ?
		// because the gin(httpRouter) does not support it, it will cause conflicts and panic directly
		comment.GET("posts/:pid/comment/:cid/children", commentController.GetPostParentCommentChildren)
		comment.POST("posts/:pid", middlewares.JwtAuth(), commentController.CreatePostComment)
	}

	favorite := r.Group("/favorites")
	{
		favoriteController := &controller.FavoriteController{}
		favorite.GET("posts/:pid", middlewares.JwtAuth(), favoriteController.CheckUserFavoritePost)
		favorite.POST("posts/:pid", middlewares.JwtAuth(), favoriteController.AddUserFavoritePost)
		favorite.DELETE("posts/:pid", middlewares.JwtAuth(), favoriteController.DeleteUserFavoritePost)
	}

	return r, nil
}
