package v1

import (
	"LearnJapan.com/internal/core/repositories"
	"LearnJapan.com/internal/entity/models"
	"LearnJapan.com/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type SessionController struct {
	//session usecase
	repo     *repositories.SessionRepo
	userRepo *repositories.UserRepo
	logger   *logger.Logger
}

func NewSessionController(r *repositories.SessionRepo, u *repositories.UserRepo, logger *logger.Logger) *SessionController {
	return &SessionController{
		repo:     r,
		userRepo: u,
		logger:   logger,
	}
}

func (u SessionController) Login(c *gin.Context) {
	user, ok := u.userRepo.FindUserByLoginAndPassword(c.PostForm("login"), c.PostForm("password"))
	if ok {
		expires, err := u.repo.Now()
		if err != nil {
			u.logger.Error(err)
		}

		expires = expires.Add(1 * time.Hour)

		newSession := models.Session{
			SessionId: u.GenerateSessionId(14),
			UserId:    user.Id,
			Expires:   expires.Format("2006-01-02 15:04:05"),
		}

		if err := u.repo.Add(&newSession); err != nil {
			u.logger.Error(err)
		}

		fmt.Println(expires)

		cookie := http.Cookie{
			Name:    "sessionId",
			Value:   newSession.SessionId,
			Expires: expires,
			Path:    "/",
		}

		http.SetCookie(c.Writer, &cookie)
		c.Redirect(302, "/")
		return
	}
}

// AuthIndex страница авторизации
func (u SessionController) AuthIndex(c *gin.Context) {
	files := []string{
		"./view/html/authorization.html",
		"./view/html/parts/header.html",
		"./view/html/parts/mainMenu.html",
	}

	templ, err := template.ParseFiles(files...)
	if err != nil {
		u.logger.Error(err)
	}

	if err := templ.Execute(c.Writer, nil); err != nil {
		u.logger.Error(err)
	}
}

// DestroySession удалит текущую сессию log out
func (u SessionController) DestroySession(c *gin.Context) {
	sessionData, _ := c.Get("sessionId")
	sessionId, _ := sessionData.(string)

	if err := u.repo.DeleteSession(sessionId); err != nil {
		u.logger.Error(err)
	}

	c.SetCookie("sessionId", "", 0, "/", "/", false, false)
	c.Redirect(http.StatusFound, "/")
}

// GenerateSessionId вернет случайный набор символов заданной длины
func (c SessionController) GenerateSessionId(len int) string {
	//TODO убрать этот метод в utils (session_usecase)
	//TODO убрать символы в конфиги
	charSet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	var sessionId []byte

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		sessionId = append(sessionId, charSet[rand.Intn(len)])
	}

	return string(sessionId)
}
