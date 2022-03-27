package main

import (
	_ "LearnJapan.com/cmd/router"
	"net/http"
)

func main(){
	http.ListenAndServe("localhost:8080", nil)
}
