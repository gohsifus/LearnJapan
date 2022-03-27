package router

import (
	"LearnJapan.com/pkg/models"
	"html/template"
	"net/http"
)

func init(){
	http.HandleFunc("/", mainIndex)
	http.HandleFunc("/dictionary/", dictionaryIndex)
}

func dictionaryIndex(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./view/html/dictionary.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil{
		panic(err)
	}

	words := models.GetById(12)

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

