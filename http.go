package main

import (
	"github.com/julienschmidt/httprouter"
    "net/http"
	"time"
	"encoding/json"
)

type Response struct {
  Status    int
  Content   string
}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := time.Now().Local()
	//fmt.Fprintf(w, t.Format("20060102150405"))
	response := Response{200, t.Format("20060102150405")}
	js, err := json.Marshal(response)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	router := httprouter.New()
	router.GET("/ping", handler)	
	
    http.ListenAndServe(":8080", router)
}