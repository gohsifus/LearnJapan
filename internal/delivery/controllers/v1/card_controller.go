package v1

import (
	"LearnJapan.com/internal/core/repositories"
	"LearnJapan.com/internal/entity/models"
	"LearnJapan.com/pkg/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

type CardController struct {
	//usecase
	cardRepo    *repositories.CardRepo
	sessionRepo *repositories.SessionRepo
	logger      *logger.Logger
}

func NewCardController(r *repositories.CardRepo, s *repositories.SessionRepo, logger *logger.Logger) *CardController {
	return &CardController{
		cardRepo:    r,
		sessionRepo: s,
		logger:      logger,
	}
}

// MainIndex Главная страница
func (u CardController) MainIndex(c *gin.Context) {
	files := []string{
		"./view/html/mainPage.html",
		"./view/html/parts/header.html",
		"./view/html/parts/footer.html",
		"./view/html/parts/mainMenu.html",
		"./view/html/parts/sitePreview.html",
	}

	data := make(map[string]interface{})

	if isReg, err := c.Cookie("isReg"); err == nil && isReg == "true" {
		data["isReg"] = "true"
	}

	if sessionId, ok := c.Get("sessionId"); ok {
		data["SessionId"] = sessionId
	}

	randCard, err := u.cardRepo.GetRandCard()
	if err != nil {
		u.logger.Error(err)
	}

	data["RandCard"] = randCard

	template, err := template.ParseFiles(files...)
	if err != nil {
		u.logger.Error(err)
	}

	if err = template.Execute(c.Writer, data); err != nil {
		u.logger.Error(err)
	}
}

// GetOneCard вернет карточку по id
func (u CardController) GetOneCard(c *gin.Context) {
	files := []string{
		"./view/html/word.html",
		"./view/html/parts/header.html",
		"./view/html/parts/footer.html",
	}

	word, err := u.cardRepo.GetCardById(c.Query("id"))
	if err != nil {
		u.logger.Error(err)
	}

	templ, err := template.ParseFiles(files...)
	if err != nil {
		u.logger.Error(err)
	}

	if err := templ.Execute(c.Writer, word); err != nil {
		u.logger.Error(err)
	}
}

// DictionaryIndex Страница словарь
func (u CardController) DictionaryIndex(c *gin.Context) {
	sessionData, _ := c.Get("sessionId")
	sessionId, _ := sessionData.(string)

	files := []string{
		"./view/html/dictionary.html",
		"./view/html/parts/header.html",
		"./view/html/parts/footer.html",
		"./view/html/parts/mainMenu.html",
		"./view/html/parts/sitePreview.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		u.logger.Error(err)
	}

	cards, err := u.cardRepo.GetCardListBySessionId(sessionId)
	if err != nil {
		u.logger.Error(err)
	}

	data := struct {
		Words     []models.JpnCard
		SessionId string
	}{
		Words:     cards,
		SessionId: sessionId,
	}

	if err := template.Execute(c.Writer, data); err != nil {
		u.logger.Error(err)
	}
}

// AddWord Добавляет карточку POST
func (u CardController) AddWord(c *gin.Context) {
	c.Request.ParseForm()

	sessionData, _ := c.Get("sessionId")
	sessionId, _ := sessionData.(string)

	userId, err := u.sessionRepo.GetUserIdBySessionId(sessionId)
	if err != nil || userId == 0 {
		u.logger.Error("user not exists")
		return
	}

	newItem := models.JpnCard{
		InJapan:   c.Request.Form.Get("InJapan"),
		InRussian: c.Request.Form.Get("InRussian"),
		Mark:      0,
		DateAdd:   c.Request.Form.Get("DateAdd"),
		UserId:    userId,
	}

	if err := u.cardRepo.Add(&newItem); err != nil {
		data := make(map[string]interface{})
		data["status"] = "Err"
		data["data"] = err

		resp, errJson := json.Marshal(data)
		if errJson != nil {
			u.logger.Error(errJson)
		}

		c.Writer.Write(resp)
	}

	data := make(map[string]interface{})
	data["status"] = "Ok"
	data["data"] = newItem

	resp, errJson := json.Marshal(data)
	if errJson != nil {
		u.logger.Error(errJson)
	}

	c.Writer.Write(resp)
}

// FindCard найдет карточку
func (u CardController) FindCard(c *gin.Context) {
	c.Request.ParseForm()

	response := make(map[string]interface{})

	if c.Request.PostForm.Get("Action") == "findByInJapan" {
		card, err := u.cardRepo.GetCardByInJapan(c.Request.PostForm.Get("InJapan"))
		if err != nil {
			u.logger.Error(err)
			return
		}

		if card.Id > 0 {
			response["status"] = "Ok"
			response["card"] = card

			json, err := json.Marshal(response)
			if err != nil {
				u.logger.Error(err)
			}

			c.Writer.Write(json)
			return
		}

		response["status"] = "Bad"
		response["card"] = "Not found"

		json, err := json.Marshal(response)
		if err != nil {
			u.logger.Error(err)
			return
		}

		c.Writer.Write(json)
	}
}

// ChangeMark изменит значение оценки слова
func (u CardController) ChangeMark(c *gin.Context) {
	c.Request.ParseForm()

	cardId, err := strconv.Atoi(c.Request.PostForm.Get("cardId"))
	if err != nil {
		u.logger.Error(err)
	}

	value, err := strconv.Atoi(c.Request.PostForm.Get("value"))
	if err != nil {
		u.logger.Error(err)
	}

	_, err = u.cardRepo.UpdateCardMark(cardId, value)
	if err != nil {
		u.logger.Error(err)
		return
	}
}

// StatisticIndex страница со статистикой по выученным словам
func (u CardController) StatisticIndex(c *gin.Context) {
	sessionData, _ := c.Get("sessionId")
	sessionId, _ := sessionData.(string)

	files := []string{
		"./view/html/statistic.html",
		"./view/html/parts/header.html",
		"./view/html/parts/mainMenu.html",
	}

	data := struct {
		BadWords  int
		NewWords  int
		AvgWords  int
		GoodWords int
		SessionId string
		AllWords  int
	}{}

	template, err := template.ParseFiles(files...)
	if err != nil {
		u.logger.Error(err)
	}

	cards, err := u.cardRepo.GetCardListBySessionId(sessionId)
	if err != nil {
		u.logger.Error(err)
	}

	for _, v := range cards {
		if v.Mark < 0 {
			data.BadWords += 1
		} else if v.Mark == 0 {
			data.NewWords += 1
		} else if v.Mark > 0 && v.Mark < 30 {
			data.AvgWords += 1
		} else if v.Mark > 30 {
			data.GoodWords += 1
		}
	}

	data.AllWords = len(cards)
	data.SessionId = sessionId

	if err := template.Execute(c.Writer, data); err != nil {
		u.logger.Error(err)
	}
}

// TestingIndex Страница с тестами
func (u CardController) TestingIndex(c *gin.Context) {
	sessionData, _ := c.Get("sessionId")
	sessionId, _ := sessionData.(string)

	files := []string{
		"./view/html/testing.html",
		"./view/html/parts/header.html",
		"./view/html/parts/mainMenu.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		u.logger.Error(err)
	}

	randCard, err := u.cardRepo.GetRandCardForUser(sessionId)
	if err != nil {
		u.logger.Error(err)
	}

	data := struct {
		RandCard  models.JpnCard
		SessionId string
	}{
		RandCard:  randCard,
		SessionId: sessionId,
	}

	if err := template.Execute(c.Writer, data); err != nil {
		u.logger.Error(err)
	}
}
