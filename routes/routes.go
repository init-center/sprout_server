package routes

import (
	"sprout_server/common/myvali"
	"sprout_server/controller"
	"sprout_server/logger"

	"github.com/gin-gonic/gin"
)

func Setup() (*gin.Engine, error) {
	r := gin.New()
	if err := myvali.Init(); err != nil {
		return r, err
	}
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// register the routes

	// create the user route group
	u := r.Group("/users")
	{
		userController := &controller.UserController{}
		u.POST("", userController.SignUp)

	}

	s := r.Group("/session")
	{
		sessionController := &controller.SessionController{}
		s.POST("", sessionController.SignIn)
	}

	c := r.Group("/vcode")
	{
		vCodeController := &controller.VCodeController{}
		c.POST("/ecode", vCodeController.SendECode)
	}

	return r, nil
}
