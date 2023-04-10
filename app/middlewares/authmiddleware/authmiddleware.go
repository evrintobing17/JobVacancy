package authmiddleware

import "github.com/gin-gonic/gin"

type AuthMiddleware interface {
	AuthorizeJWTWithUserContext() gin.HandlerFunc
}
