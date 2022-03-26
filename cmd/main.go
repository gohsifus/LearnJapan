package main

import (
	_ "./router"
	"net/http"
)

func main(){
	http.ListenAndServe("localhost:8080", nil)
}
