package router

import (
	"LearnJapan.com/pkg/models"
	"html/template"
	"net/http"
)

func init(){
	http.HandleFunc("/", mainIndex)
	http.HandleFunc("/dictionary/", dictionaryIndex)
	http.HandleFunc("/dictionary/oneWord/", getOneCard)
}

func dictionaryIndex(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/dictionary.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil{
		panic(err)
	}

	words := models.GetList()

	err = template.Execute(w, words)
	if err != nil{
		panic(err)
	}
}

func mainIndex(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/mainPage.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil{
		panic(err)
	}

	err = template.Execute(w, nil)
	if err != nil{

		panic(err)
	}
}

func getOneCard(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/word.html",
	}

	word := models.GetById(r.URL.Query().Get("id"))
	templ, err := template.ParseFiles(files...)
	if err != nil{
		panic(err)
	}

	templ.Execute(w, word)
}

