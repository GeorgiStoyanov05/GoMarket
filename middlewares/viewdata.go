package middlewares

import "github.com/gin-gonic/gin"

func WithAuth(c *gin.Context, data gin.H) gin.H {
	if data == nil {
		data = gin.H{}
	}

	// default is false if middleware didn't set it
	if v, ok := c.Get("IsLoggedIn"); ok {
		data["IsLoggedIn"] = v
	} else {
		data["IsLoggedIn"] = false
	}

	if u, ok := c.Get("user"); ok {
		data["user"] = u
	}

	return data
}
