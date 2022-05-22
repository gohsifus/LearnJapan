package router

import (
	"LearnJapan.com/pkg/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func init(){
	http.HandleFunc("/", mainIndex)
	http.HandleFunc("/dictionary/", dictionaryIndex)
	http.HandleFunc("/dictionary/oneWord/", getOneCard)
	http.HandleFunc("/dictionary/addWord/", addWord)
	http.HandleFunc("/registration/", registrationIndex)
	http.HandleFunc("/registration/addUser/", addUser)

	//Обработка статических файлов
	fileServer := http.FileServer(http.Dir("./view/static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
}

//dictionaryIndex Страница словарь
func dictionaryIndex(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/dictionary.html",
		"./view/html/parts/header.html",
		"./view/html/parts/footer.html",
		"./view/html/parts/mainMenu.html",
		"./view/html/parts/sitePreview.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil{
		panic(err)
	}

	words := struct{
		Words []models.JpnCards
	}{
		Words: models.GetCardList(),
	}

	err = template.Execute(w, words)
	if err != nil{
		panic(err)
	}
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

	template, err := template.ParseFiles(files...)
	if err != nil{
		panic(err)
	}

	data := make(map[string]string)
	if isReg := r.URL.Query().Get("isReg"); isReg == "true" {
		data = map[string]string{
			"isReg": "true",
		}
	}

	err = template.Execute(w, data)
	if err != nil{

		panic(err)
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
		panic(err)
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
		panic(err)
	}

	templ.Execute(w, nil)
}

//addWord Добавляет карточку POST
func addWord(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		r.ParseForm()

		newItem := models.JpnCards{
			InJapan: r.Form.Get("InJapan"),
			InRussian: r.Form.Get("InRussian"),
			Mark: 0,
			DateAdd: r.Form.Get("DateAdd"),
		}

		if ok, err := newItem.Add(); err == nil && ok{
			fmt.Println("Запись добавлена")

			data := make(map[string]interface{})
			data["status"] = "Ok"
			data["data"] = newItem

			resp, errJson := json.Marshal(data)
			if errJson != nil{
				panic(errJson)
			}
			w.Write(resp)
		} else {
			fmt.Println("Запись не добавлена: ", err)

			data := make(map[string]interface{})
			data["status"] = "Err"
			data["data"] = err

			resp, errJson := json.Marshal(data)
			if errJson != nil{
				panic(errJson)
			}

			w.Write(resp)
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
			fmt.Println("Ok")
			http.Redirect(w, r, "/?isReg=true", http.StatusSeeOther)
		} else {
			fmt.Println("Error: ", err)
		}
	} else {
		//TODO 404 сделать страницу
	}
}

