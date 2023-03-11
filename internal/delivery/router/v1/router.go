package v1

import (
	"LearnJapan.com/constants"
	"LearnJapan.com/internal/delivery/controllers"
	"LearnJapan.com/internal/delivery/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	Mux            *gin.Engine
	Controllers    *controllers.MainController
	AuthMiddleware *middlewares.AuthMiddleware
}

func NewRouter(m *gin.Engine, c *controllers.MainController, a *middlewares.AuthMiddleware) *Router {
	return &Router{
		Mux:            m,
		Controllers:    c,
		AuthMiddleware: a,
	}
}

func (r Router) Setup() {
	//r.Mux.Use(gin.Logger())

	dictionary := r.Mux.Group(constants.ROUTE_DICTIONARY)
	dictionary.Use(r.AuthMiddleware.Access())

	testing := r.Mux.Group(constants.ROUTE_TESTING)
	testing.Use(r.AuthMiddleware.Access())

	registration := r.Mux.Group(constants.ROUTE_REGISTRATION)
	authorization := r.Mux.Group(constants.ROUTE_AUTHORIZATION)

	r.Mux.GET(constants.ROUTE_INDEX, r.Controllers.Cards.MainIndex)

	testing.GET(constants.ROUTE_INDEX, r.Controllers.Cards.TestingIndex)

	dictionary.GET(constants.ROUTE_INDEX, r.Controllers.Cards.DictionaryIndex)
	dictionary.POST(constants.ROUTE_ADD_CARD, r.Controllers.Cards.AddWord)
	dictionary.POST(constants.ROUTE_CHANGE_MARK, r.Controllers.Cards.ChangeMark)
	dictionary.POST(constants.ROUTE_FIND, r.Controllers.Cards.FindCard)
	dictionary.GET(constants.ROUTE_STATISTIC, r.Controllers.Cards.StatisticIndex)

	registration.GET(constants.ROUTE_INDEX, r.Controllers.RegistrationIndex)
	registration.POST(constants.ROUTE_ADD_USER, r.Controllers.Users.AddUser)

	authorization.GET(constants.ROUTE_INDEX, r.Controllers.Sessions.AuthIndex)
	authorization.POST(constants.ROUTE_INDEX, r.Controllers.Sessions.Login)
	authorization.GET(constants.ROUTE_LOGOUT, r.Controllers.Sessions.DestroySession)

	//Обработка статических файлов
	r.Mux.StaticFS("/static", http.Dir("./view/static"))
}
