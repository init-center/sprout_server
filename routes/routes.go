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
		AllowOrigins: []string{
			"https://init.center",
			"http://init.center",
			"https://admin.init.center",
			"http://admin.init.center",
			"https://blog.init.center",
			"http://admin.blog.init.center",
			"https://admin.blog.init.center",
			"http://blog.init.center",
			"http://localhost:3000",
			"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"content-type", "Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}), middlewares.ParseOrigin())

	// register the routes

	// create the user route group
	user := r.Group("/users")
	{
		userController := &controller.UserController{}
		user.POST("", userController.SignUp)
		user.GET("/public/:uid", userController.GetPublicUserInfo)
		user.GET("/private", middlewares.JwtAuth(), userController.GetPrivateUserInfo)
		user.PUT("", middlewares.JwtAuth(), userController.UpdateUser)
		user.DELETE("", middlewares.JwtAuth(), userController.DeleteUser)
		user.PUT("/password", userController.UpdatePassword)
		user.GET("/ban_time/:uid", userController.GetBanTime)

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

	configR := r.Group("/configs")
	{
		configController := &controller.ConfigController{}
		configR.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), configController.Create)
		configR.PUT("/:key", middlewares.JwtAuth(), middlewares.AdminAuth(), configController.Update)
		configR.DELETE("/:key", middlewares.JwtAuth(), middlewares.AdminAuth(), configController.Delete)
		configR.GET("", configController.GetByQuery)
		configR.GET("/:key", configController.GetConfigByKey)
	}

	post := r.Group("/posts")
	{
		postController := &controller.PostController{}
		post.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), postController.Create)
		post.PUT("/:pid", middlewares.JwtAuth(), middlewares.AdminAuth(), postController.Update)
		post.GET("", postController.GetPostList)
		post.GET("/detail", postController.GetPostDetailList)
		post.GET("/detail/:pid", postController.GetPostDetail)
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

		adminAnalysis := admin.Group("/analysis")
		{
			adminAnalysisController := &controller.AnalysisController{}
			adminAnalysis.GET("/users", adminAnalysisController.GetUserAnalysis)
			adminAnalysis.GET("/comments", adminAnalysisController.GetCommentAnalysis)
			adminAnalysis.GET("/views", adminAnalysisController.GetViewsAnalysis)
			adminAnalysis.GET("/posts", adminAnalysisController.GetPostAnalysis)
			adminAnalysis.GET("/complex", adminAnalysisController.GetComplexAnalysis)
			adminAnalysis.GET("/posts/viewsrank", adminAnalysisController.GetPostViewsRank)
			adminAnalysis.GET("/categories/postscount", adminAnalysisController.GetCategoriesPostsCount)
			adminAnalysis.GET("/tags/postscount", adminAnalysisController.GetTagsPostsCount)
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
		comment.GET("", commentController.GetPublicComments)
	}

	favorite := r.Group("/favorites")
	{
		favoriteController := &controller.FavoriteController{}
		favorite.GET("/posts/:pid", middlewares.JwtAuth(), favoriteController.CheckUserFavoritePost)
		favorite.POST("/posts/:pid", middlewares.JwtAuth(), favoriteController.AddUserFavoritePost)
		favorite.DELETE("/posts/:pid", middlewares.JwtAuth(), favoriteController.DeleteUserFavoritePost)
		favorite.GET("", favoriteController.GetByQuery)
	}

	pageViews := r.Group("/views")
	{
		pageViewsController := &controller.PageViewsController{}
		pageViews.POST("", pageViewsController.CreatePageViews)
		pageViews.GET("", pageViewsController.GetPageViews)
	}

	friend := r.Group("/friends")
	{
		friendController := &controller.FriendController{}
		friend.POST("", middlewares.JwtAuth(), middlewares.AdminAuth(), friendController.Create)
		friend.GET("", friendController.GetByQuery)
		friend.PUT("/:id", middlewares.JwtAuth(), middlewares.AdminAuth(), friendController.Update)
		friend.DELETE("/:id", middlewares.JwtAuth(), middlewares.AdminAuth(), friendController.Delete)
	}

	return r, nil
}
