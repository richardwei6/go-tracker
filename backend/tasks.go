package main

import (
	"os"
	"log"
	"net/http"
)

const(
	JSONFilePath = "./tasks.json"
)

type Task struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Done bool `json:"done"`
}

type CreateTaskBody struct {
	Name string `json:"name"`

}

func getTasks(w http.ResponseWriter, _ *http.Request){ // returns saved json
	json, err := os.ReadFile("tasks.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to read json file %v", err)
		return
	}

	w.Write(json)
}

/*func createTask(w http.ResponseWriter, _ *http.Request){ // creates new task based on request

}*/