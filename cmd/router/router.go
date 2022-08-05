package router

import (
	"LearnJapan.com/pkg/logger"
	"LearnJapan.com/pkg/models"
	"LearnJapan.com/pkg/yandexTranslateApi"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func init(){
	http.HandleFunc("/", mainIndex)
	http.HandleFunc("/testing/", testingIndex)
	http.HandleFunc("/dictionary/", dictionaryIndex)
	http.HandleFunc("/dictionary/oneWord/", getOneCard)
	http.HandleFunc("/dictionary/addWord/", addWord)
	http.HandleFunc("/dictionary/changeMark/", changeMark)
	http.HandleFunc("/registration/", registrationIndex)
	http.HandleFunc("/registration/addUser/", addUser)
	http.HandleFunc("/authorization/", authIndex)
	http.HandleFunc("/authorization/exit/", destroySession)
	http.HandleFunc("/dictionary/find/", findCard)
	http.HandleFunc("/dictionary/translate", translate)
	http.HandleFunc("/statistic/", statisticIndex)

	//Обработка статических файлов
	fileServer := http.FileServer(http.Dir("./view/static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
}

//mainIndex Главная страница
func mainIndex(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/mainPage.html",
		"./view/html/parts/header.html",
		"./view/html/parts/footer.html",
		"./view/html/parts/mainMenu.html",
		"./view/html/parts/sitePreview.html",
	}

	data := make(map[string]interface{})

	if isReg, err := r.Cookie("isReg"); err == nil && isReg.Value == "true"{
		data["isReg"] = "true"
	}

	if sessionId, ok := getCurrentSessionId(r); ok{
		data["SessionId"] = sessionId
	}

	randCard, err :=  models.GetRandCard()
	if err != nil{
		logger.Print("Error: RandCard" + err.Error())
	}
	data["RandCard"] = randCard

	template, err := template.ParseFiles(files...)
	if err != nil{
		logger.Print("Error: templateParse" + err.Error())
	}

	err = template.Execute(w, data)
	if err != nil{
		panic(err)
		//TODO Возвращать код ошибки сервера
	}
}

//testingIndex Страница с тестами
func testingIndex(w http.ResponseWriter, r *http.Request){
	if checkAccess(r){
		files := []string{
			"./view/html/testing.html",
			"./view/html/parts/header.html",
			"./view/html/parts/mainMenu.html",
		}

		template, err := template.ParseFiles(files...)
		if err != nil {
			logger.Print("Error: templateParse " + err.Error())
		}

		sessionId, _ := r.Cookie("sessionId")
		randCard, err := models.GetRandCardForUser(sessionId.Value)
		if err != nil {
			logger.Print(err.Error())
		}


		data := struct{
			RandCard models.JpnCards
			SessionId string
		}{
			RandCard: randCard,
			SessionId: sessionId.Value,
		}

		err = template.Execute(w, data)
		if err != nil{
			panic(err)
		}
	} else {
		http.Redirect(w, r, "/authorization/", 302)
	}
}

//dictionaryIndex Страница словарь
func dictionaryIndex(w http.ResponseWriter, r *http.Request){
	if checkAccess(r) {
		files := []string{
			"./view/html/dictionary.html",
			"./view/html/parts/header.html",
			"./view/html/parts/footer.html",
			"./view/html/parts/mainMenu.html",
			"./view/html/parts/sitePreview.html",
		}

		template, err := template.ParseFiles(files...)
		if err != nil {
			logger.Print("Error: templateParse " + err.Error())
		}

		sessionId, _ := r.Cookie("sessionId")
		cards, err := models.GetCardListBySessionId(sessionId.Value)
		if err != nil{
			logger.Print("Error: GetCardList " + err.Error())
		}

		data := struct {
			Words []models.JpnCards
			SessionId string
		}{
			Words: cards,
			SessionId: sessionId.Value,
		}

		err = template.Execute(w, data)
		if err != nil {
			panic(err)
		}
	} else {
		http.Redirect(w, r, "/authorization/", 302)
	}
}

//getOneCard Вернет карточку по id Get
func getOneCard(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/word.html",
		"./view/html/parts/header.html",
		"./view/html/parts/footer.html",
	}

	word := models.GetCardById(r.URL.Query().Get("id"))
	templ, err := template.ParseFiles(files...)
	if err != nil{
		logger.Print("Error: templateParse " + err.Error())
	}

	templ.Execute(w, word)
}

//registrationIndex Страница регистрации пользователя
func registrationIndex(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/registration.html",
		"./view/html/parts/header.html",
		"./view/html/parts/mainMenu.html",
	}

	templ, err := template.ParseFiles(files...)
	if err != nil{
		logger.Print("Error: templateParse " + err.Error())
	}

	templ.Execute(w, nil)
}

//addWord Добавляет карточку POST
func addWord(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		r.ParseForm()

		sessionId, err := r.Cookie("sessionId")
		if err != nil{
			logger.Print(err.Error())
		}
		userId, ok := models.GetUserIdBySessionId(sessionId.Value)
		if !ok {
			//TODO Перенаправление на авторизацию
			fmt.Println("Недействительная сессия или ошибка запроса")
		} else {
			newItem := models.JpnCards{
				InJapan:   r.Form.Get("InJapan"),
				InRussian: r.Form.Get("InRussian"),
				Mark:      0,
				DateAdd:   r.Form.Get("DateAdd"),
				UserId: userId,
			}

			if ok, err := newItem.Add(); err == nil && ok {

				data := make(map[string]interface{})
				data["status"] = "Ok"
				data["data"] = newItem

				resp, errJson := json.Marshal(data)
				if errJson != nil {
					logger.Print(err.Error())
				}
				w.Write(resp)
			} else {

				data := make(map[string]interface{})
				data["status"] = "Err"
				data["data"] = err

				resp, errJson := json.Marshal(data)
				if errJson != nil {
					logger.Print(err.Error())
				}

				w.Write(resp)
			}
		}
	} else {
		//TODO 404 сделать страницу
	}
}

//addUser Зарегистрирует нового пользователя POST
func addUser(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()

		newUser := models.User{
			Login: r.Form.Get("login"),
			Password: r.Form.Get("password"),
			Email: r.Form.Get("email"),
		}

		if Ok, err := newUser.Add();  Ok{
			//Передаем куку чтобы вывести информацию о регистрации на главной странице
			expiration := models.Now().Add(2 * time.Second)
			cookie := http.Cookie{Name: "isReg", Value: "true", Expires: expiration, Path: "/"}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			logger.Print(err.Error())
		}
	} else {
		//TODO 404 сделать страницу
	}
}

//authIndex Страница авторизации
func authIndex(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" { //Авторизация
		r.ParseForm()

		user, ok := models.FindUserByLoginAndPassword(r.PostForm.Get("login"), r.PostForm.Get("password"))
		if ok{
			fmt.Println(models.Now())
			expires := models.Now().Add(1 * time.Hour)
			newSession := models.Session{
				SessionId: generateSessionId(14),
				UserId: user.Id,
				Expires: expires.Format("2006-01-02 15:04:05"),
			}

			newSession.Add()

			cookie := http.Cookie{
				Name: "sessionId",
				Value: newSession.SessionId,
				Expires: expires,
				Path: "/",
			}

			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", 302)
		}
	}

	files := []string{
		"./view/html/authorization.html",
		"./view/html/parts/header.html",
		"./view/html/parts/mainMenu.html",
	}

	templ, err := template.ParseFiles(files...)
	if err != nil {
		logger.Print(err.Error())
	}

	templ.Execute(w, nil)
}

//generateSessionId Вернет случайный набор символов заданной длины
func generateSessionId(len int) string{
	charSet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	var sessionId []byte

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++{
		sessionId = append(sessionId, charSet[rand.Intn(len)])
	}

	return string(sessionId)
}

//checkAccess Проверит есть ли доступ к странице
func checkAccess(r *http.Request) bool{
	if cookieVal, err := r.Cookie("sessionId"); err == nil{
		if ok, errAlive := models.IsAliveSession(cookieVal.Value); ok && errAlive == nil{
			return true
		} else {
			logger.Print("Ошибка проверки сессии " + errAlive.Error())
			return false
		}
	} else if err != nil{
		logger.Print("Ошибка авторизации " + err.Error())
	}

	return false
}

//getCurrentSessionId Вернет текущую сессию, если она действует
func getCurrentSessionId(r *http.Request) (string, bool){
	if checkAccess(r){
		currentSession, _ := r.Cookie("sessionId")
		return currentSession.Value, true
	}

	return "", false
}

//destroySession Удалит текущую сессию log out
func destroySession(w http.ResponseWriter, r *http.Request){
	sessionID, _ := r.Cookie("sessionId")

	models.DeleteSession(sessionID.Value)
	http.Redirect(w, r, "/", 302)
}

//findCard Найдет карточку
func findCard(w http.ResponseWriter, r * http.Request){
	if r.Method == "POST"{
		r.ParseForm()

		response := make(map[string]interface{})

		if r.PostForm.Get("Action") == "findByInJapan"{
			card, ok := models.GetCardByInJapan(r.PostForm.Get("InJapan"))
			if ok {
				response["status"] = "Ok"
				response["card"] = card

				json, err := json.Marshal(response)
				if err != nil{
					panic(err)
				}

				w.Write(json)
			} else {
				response["status"] = "Bad"
				response["card"] = "Not found"

				json, err := json.Marshal(response)
				if err != nil{
					logger.Print(err.Error())
				}

				w.Write(json)
			}
		}
	}
}

//translate Переведет слово с помощью api
func translate(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()

		response := make(map[string]string)

		var srcStr string
		if r.PostForm.Get("inJapan") != ""{
			srcStr = r.PostForm.Get("inJapan")
		} else {
			srcStr = r.PostForm.Get("inRussia")
		}

		translatedString, err := yandexTranslateApi.Translate(
															srcStr,
															r.PostForm.Get("srcCode"),
															r.PostForm.Get("dstCode"))
		if err != nil{
			response = map[string]string{
				"status": "err",
				"explain": err.Error(),
			}
		} else {
			response = map[string]string{
				"status": "ok",
				"translated": translatedString,
			}
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil{
			logger.Print("Translate Error " + err.Error())
		}
		w.Write(jsonResponse)
	}
}

//changeMark
func changeMark(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST"{
		r.ParseForm()

		if ok, _ := models.IsAliveSession(r.PostForm.Get("sessionId")); ok {
			cardId, err := strconv.Atoi(r.PostForm.Get("cardId"))
			if err != nil {
				fmt.Println(err)
			}

			value, err := strconv.Atoi(r.PostForm.Get("value"))
			if err != nil{
				fmt.Println(err)
			}

			err = models.UpdateCardMark(cardId, value)
			if err != nil{
				fmt.Println(err)
			}
		}
	}
}

//statisticIndex Страница с статистикой
func statisticIndex(w http.ResponseWriter, r *http.Request){
	if checkAccess(r) {
		files := []string{
			"./view/html/statistic.html",
			"./view/html/parts/header.html",
			"./view/html/parts/mainMenu.html",
		}

		data := struct{
			BadWords int
			NewWords int
			AvgWords int
			GoodWords int
			SessionId string
			AllWords int
		}{}

		template, err := template.ParseFiles(files...)
		if err != nil {
			logger.Print("Error: templateParse " + err.Error())
		}

		sessionId, _ := r.Cookie("sessionId")
		cards, err := models.GetCardListBySessionId(sessionId.Value)
		if err != nil{
			logger.Print("Error: GetCardList " + err.Error())
		}

		for _, v := range cards{
			if v.Mark < 0{
				data.BadWords += 1
			} else if v.Mark == 0{
				data.NewWords += 1
			} else if v.Mark > 0 && v.Mark < 30{
				data.AvgWords += 1
			} else if v.Mark > 30{
				data.GoodWords += 1
			}
		}

		data.AllWords = len(cards)
		data.SessionId = sessionId.Value

		err = template.Execute(w, data)
		if err != nil {
			panic(err)
		}
	} else {
		http.Redirect(w, r, "/authorization/", 302)
	}
}

