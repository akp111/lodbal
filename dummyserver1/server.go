package main

import "net/http"
import "log"

func main(){
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health is top notch")
	})
	log.Println("Server started")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Server crashed")
	}
}