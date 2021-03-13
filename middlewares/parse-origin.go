package middlewares

import (
	"sprout_server/common/constants"

	"github.com/mssola/user_agent"

	"github.com/gin-gonic/gin"
)

func ParseOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("user-agent")
		ua := user_agent.New(userAgent)
		os := ua.OS()
		engineName, engineVersion := ua.Engine()
		engine := engineName
		if engineVersion != "" {
			engine += " " + engineVersion
		}

		browserName, browserVersion := ua.Browser()
		browser := browserName
		if browserVersion != "" {
			browser += " " + browserVersion
		}

		ip := c.ClientIP()
		// set the info to gin context
		c.Set(constants.CtxOriginIpKey, ip)
		c.Set(constants.CtxOriginUAKey, userAgent)
		c.Set(constants.CtxOriginOsKey, os)
		c.Set(constants.CtxOriginEngineKey, engine)
		c.Set(constants.CtxOriginBrowserKey, browser)
		c.Next()
	}
}
