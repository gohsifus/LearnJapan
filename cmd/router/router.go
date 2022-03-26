package router

import (
	"html/template"
	"net/http"
)

func init(){
	http.HandleFunc("/", mainIndex)
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

