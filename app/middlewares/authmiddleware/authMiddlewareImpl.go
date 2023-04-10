package authmiddleware

import (
	"errors"
	"strings"

	"github.com/evrintobing17/JobVacancy/app/helpers/jsonhttpresponse"
	"github.com/evrintobing17/JobVacancy/app/helpers/jwthelper"
	job "github.com/evrintobing17/JobVacancy/app/modules/job"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidToken          = errors.New("invalid token")
	ErrUserContextNotSet     = errors.New("user context is empty. Use AuthorizeJWTWithUserContext instead")
	ErrInvalidResourceAccess = errors.New("this user has no rights to access this resource")
)

type authMiddleware struct {
	jobService job.JobRepository
}

func NewAuthMiddleware(jobService job.JobRepository) AuthMiddleware {
	return &authMiddleware{jobService: jobService}
}

//AuthorizeJWTWithUserContext - Authorize JWT with User Context (Need to look up for user in DB in every request)
func (auth *authMiddleware) AuthorizeJWTWithUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		//Get User Claims
		if bearerToken == "" {
			jsonhttpresponse.Unauthorized(c, ErrInvalidToken.Error())
			c.Abort()
			return
		}

		//Extract JWT Token from Bearer
		jwtTokenSplit := strings.Split(bearerToken, "Bearer ")
		if jwtTokenSplit[1] == "" {
			jsonhttpresponse.Unauthorized(c, ErrInvalidToken.Error())
			c.Abort()
			return
		}
		jwtToken := jwtTokenSplit[1]

		jwtTokenClaims, err := jwthelper.VerifyTokenWithClaims(jwtToken)
		if err != nil {
			jsonhttpresponse.Unauthorized(c, ErrInvalidToken.Error())
			c.Abort()
			return
		}

		if jwtTokenClaims.Valid() != nil {
			jsonhttpresponse.Unauthorized(c, ErrInvalidToken.Error())
		}

		c.Next()
		return
	}
}
