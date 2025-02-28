package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	json, err := os.ReadFile(JSONFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to read json file %v", err)
		return
	}

	w.Write(json)
}

func createTask(w http.ResponseWriter, r *http.Request){ // creates new task based on request
	body := CreateTaskBody{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil { // decode http body json
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to decode json file - %v", err)
		return
	}
	
	jsonFile, err := os.Open(JSONFilePath)
	if err != nil { // open local json file
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to open json file - %v", err)
		return
	}

	tasks := []Task{}
	if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil { // decode local json file
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to decode local json file into arr %v", err)
		return
	}

	///

	newTask  := Task {
		ID: len(tasks) + 1,
		Name: body.Name,
		Done: false,
	}

	tasks = append(tasks, newTask)

	///

	j, err := json.Marshal(tasks)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to marshalling tasks %v", err)
		return
	}

	err = os.WriteFile(JSONFilePath, j, 0755)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error writing new json file %v", err)
		return
	}

	///

	j, err = json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshalling NEW task %v", err)
		return
	}

	// update header
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
