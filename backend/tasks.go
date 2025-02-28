package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/go-chi/chi/v5"
)

const(
	JSONFilePath = "./tasks.json"
)

var nextTaskID int // keeps track of next unused ID

type Task struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Done bool `json:"done"`
}

type CreateTaskBody struct {
	Name string `json:"name"`
}

type UpdateTaskBody struct { // use * to check if non nil
	Name *string `json:"name"`
	Done *bool `json:"done"`
}

///

func loadNextTaskID() {
	tasks, err := readLocalJSON()
	if err != nil {
		log.Printf("didn't find a previous JSON file")
		return
	} 

	var mx int = -1
	for _, task := range tasks {
		mx = max(mx, task.ID)
	}

	if mx == -1 {
		log.Printf("didn't find any tasks?")
		return
	}

	nextTaskID = mx + 1
	log.Printf("next task ID = %v", nextTaskID)
}

//

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
		log.Printf("Unable to decode request - %v", err)
		return
	}
	
	tasks, err := readLocalJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to load local json file - %v", err)
		return
	}

	///

	newTask := Task {
		ID: nextTaskID,
		Name: body.Name,
		Done: false,
	}

	tasks = append(tasks, newTask)
	nextTaskID += 1

	///

	err = writeLocalJSON(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error writing new task - %v", err)
		return
	}

	///

	j, err := json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshalling NEW task %v", err)
		return
	}

	// update header
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func updateTask(w http.ResponseWriter, r *http.Request) { // update task
	body := UpdateTaskBody{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil { // decode http body json
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to decode request - %v", err)
		return
	}
	
	tasks, err := readLocalJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to load local json file - %v", err)
		return
	}

	taskRequestID, err := strconv.Atoi(chi.URLParam(r, "taskID")) // obtain taskID from URL
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to convert taskID from URL - %v", err)
		return
	}

	///

	for i, task := range tasks {
		if task.ID == taskRequestID{

			if body.Name != nil {
				task.Name = *body.Name
			}

			if body.Done != nil {
				task.Done = *body.Done
			}

			tasks[i] = task // update with new copy in array
			break
		}
	}

	///

	err = writeLocalJSON(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error writing updated task - %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) { // delete task
	tasks, err := readLocalJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to load local json file - %v", err)
		return
	}

	taskRequestID, err := strconv.Atoi(chi.URLParam(r, "taskID")) // obtain taskID from URL
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Unable to convert taskID from URL - %v", err)
		return
	}

	///
	f := false

	for i, task := range tasks {
		if task.ID == taskRequestID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			f = true
			break
		}
	}

	if (!f) { // task was not found
		w.WriteHeader(http.StatusNotFound) 
		return
	}

	///

	err = writeLocalJSON(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error writing new json with deleted entry - %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}


///

func readLocalJSON() ([]Task, error) {
	tasks := []Task{}

	jsonFile, err := os.Open(JSONFilePath)
	if err != nil { // open local json file
		return tasks, err
	}

	if err := json.NewDecoder(jsonFile).Decode(&tasks); err != nil { // decode local json file
		return tasks, err
	}

	return tasks, nil
}

func writeLocalJSON(tasks []Task) (error) {
	j, err := json.Marshal(tasks)

	if err != nil {
		return err
	}

	err = os.WriteFile(JSONFilePath, j, 0755)
	if err != nil {
		return err
	}

	return nil
}