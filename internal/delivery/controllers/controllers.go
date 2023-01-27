package controllers

import (
	"LearnJapan.com/internal/core/repositories"
	v1 "LearnJapan.com/internal/delivery/controllers/v1"
	"LearnJapan.com/pkg/logger"
	"database/sql"
	"github.com/gin-gonic/gin"
	"html/template"
)

type MainController struct {
	Cards    *v1.CardController
	Sessions *v1.SessionController
	Users    *v1.UserController
	logger   *logger.Logger
}

func NewMainController(db *sql.DB, logger *logger.Logger) *MainController {
	cardRepo := repositories.NewCardRepo(db)
	sessionRepo := repositories.NewSessionRepo(db)
	userRepo := repositories.NewUserRepo(db)

	return &MainController{
		Cards:    v1.NewCardController(cardRepo, sessionRepo, logger),
		Sessions: v1.NewSessionController(sessionRepo, userRepo, logger),
		Users:    v1.NewUserController(userRepo, sessionRepo, logger),
		logger:   logger,
	}
}

// RegistrationIndex страница регистрации пользователя
func (u MainController) RegistrationIndex(c *gin.Context) {
	files := []string{
		"./view/html/registration.html",
		"./view/html/parts/header.html",
		"./view/html/parts/mainMenu.html",
	}

	templ, err := template.ParseFiles(files...)
	if err != nil {
		u.logger.Error(err)
	}

	templ.Execute(c.Writer, nil)
}

// translate Переведет слово с помощью api
/*func (c MainController) Translate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		response := make(map[string]string)

		var srcStr string
		if r.PostForm.Get("inJapan") != "" {
			srcStr = r.PostForm.Get("inJapan")
		} else {
			srcStr = r.PostForm.Get("inRussia")
		}

		translatedString, err := yandexTranslateApi.Translate(
			srcStr,
			r.PostForm.Get("srcCode"),
			r.PostForm.Get("dstCode"))
		if err != nil {
			response = map[string]string{
				"status":  "err",
				"explain": err.Error(),
			}
		} else {
			response = map[string]string{
				"status":     "ok",
				"translated": translatedString,
			}
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Print("Translate Error " + err.Error())
		}
		w.Write(jsonResponse)
	}
}*/
