package v1

import (
	"LearnJapan.com/constants"
	"LearnJapan.com/internal/core/repositories"
	"LearnJapan.com/internal/entity/models"
	"LearnJapan.com/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserController struct {
	//usecase
	userRepo    *repositories.UserRepo
	sessionRepo *repositories.SessionRepo
	logger      *logger.Logger
}

func NewUserController(r *repositories.UserRepo, s *repositories.SessionRepo, logger *logger.Logger) *UserController {
	return &UserController{
		userRepo:    r,
		sessionRepo: s,
		logger:      logger,
	}
}

// AddUser зарегистрирует нового пользователя POST
func (u UserController) AddUser(c *gin.Context) {
	c.Request.ParseForm()

	newUser := models.User{
		Login:    c.Request.Form.Get("login"),
		Password: c.Request.Form.Get("password"),
		Email:    c.Request.Form.Get("email"),
	}

	if err := u.userRepo.Add(&newUser); err != nil {
		u.logger.Error(err)
		return
	}

	//Передаем куку, чтобы вывести информацию о регистрации на главной странице
	expiration, err := u.sessionRepo.Now()
	if err != nil {
		u.logger.Error(err)
		return
	}

	expiration = expiration.Add(2 * time.Second)

	cookie := http.Cookie{Name: "isReg", Value: "true", Expires: expiration, Path: "/"}

	http.SetCookie(c.Writer, &cookie)
	c.Redirect(http.StatusFound, constants.ROUTE_INDEX)
}
