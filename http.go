package main

import (
	"github.com/julienschmidt/httprouter"
    "net/http"
//	"time"
	"encoding/json"
//	"strconv"
//	"log"
//	"io/ioutil"
//	"net/http/httputil"
//	"fmt"
	"os"
	"encoding/gob"
)

//var log = logging.MustGetLogger("example")

var storage map[string]string

type KeyItem struct {
	Key 		string `json:"key"`
	Content 	string `json:"content"`
}

type Response struct {
  Status    int
  Content   string
}

func addKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var rec KeyItem
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil || rec.Key == "" || rec.Content == "" {
		w.WriteHeader(400)
		return
	}
	storage[rec.Key] = rec.Content

	js, err := json.Marshal(rec)
	w.WriteHeader(201) 
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func getKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	rec := KeyItem{Key: key}
	if content, ok := storage[key]; ok {
		rec.Content = content
	} else {
		w.WriteHeader(404)
		return
	}	

	js, err := json.Marshal(rec)
	w.WriteHeader(200) 
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func deleteKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	if _, ok := storage[key]; ok {
		delete(storage, key)
	} else {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200) 
}	
	

func persistToFile(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func extractFromFile(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}


func main() {
	storage = make(map[string]string)
	router := httprouter.New()
	router.GET("/keys/:key", getKey)
	router.POST("/keys", addKey)
	router.DELETE("/keys/:key", deleteKey)
    http.ListenAndServe(":8080", router)
}