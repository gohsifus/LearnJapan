package middlewares

import (
	"LearnJapan.com/constants"
	"LearnJapan.com/internal/core/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AuthMiddleware struct {
	//session usecase
	repo *repositories.SessionRepo
}

func NewAuthMiddleware(r *repositories.SessionRepo) *AuthMiddleware {
	return &AuthMiddleware{
		repo: r,
	}
}

// Access доступ к ресурсу
func (m AuthMiddleware) Access() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentSession, err := c.Cookie("sessionId")
		if err != nil {
			c.Redirect(http.StatusFound, constants.ROUTE_AUTHORIZATION)
			c.Abort()
		}

		ok, err := m.repo.IsAliveSession(currentSession)
		if err != nil {
			log.Fatal(err)
		}

		if ok {
			c.Set("sessionId", currentSession)
			c.Next()
		} else {
			c.Redirect(http.StatusFound, constants.ROUTE_AUTHORIZATION)
			c.Abort()
		}
	}
}
