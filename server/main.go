package main

import (
	"encoding/json"
	"faasedge-dag/server/dag"
	"faasedge-dag/server/scheduler"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type AddDagResponse struct {
	Message string `json:"message"`
}

var dagList = make(map[string]*dag.Dag)

func main() {
	http.HandleFunc("/dag/add", addDagHandler)
	http.HandleFunc("/dag/schedule", scheduleDagHandler)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func scheduleDagHandler(w http.ResponseWriter, r *http.Request) {
	clientId := r.Header.Get("Client-Id")
	if clientId == "" {
		http.Error(w, "Client-Id header is missing", http.StatusBadRequest)
		return
	}

	queryParameters := r.URL.Query()
	dagName := queryParameters.Get("dagName")

	if dagName == "" {
		log.Println("dagName query parameter is missing")
		http.Error(w, "dagName query parameter is missing", http.StatusBadRequest)
		return
	}

	dag, ok := dagList[dagName]
	
	if !ok {
			http.Error(w, "DAG ("+dagName+") not found", http.StatusNotFound)
			return
	}

	vdag := &scheduler.VDag{
			ClientId:      clientId,
			DagDefinition: *dag,
			PDagMap:       make(map[string][]*scheduler.PDagDeployment),
	}
	
	scheduleResponse := scheduler.Scheduler.ScheduleDag(vdag)
	
	yamlData, err := yaml.Marshal(scheduleResponse)
	if err != nil {
			http.Error(w, "Failed to marshal YAML", http.StatusInternalServerError)
			return
	}
	
	w.Header().Set("Content-Type", "application/x-yaml")
	
	_, err = w.Write(yamlData)
	if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
	}
	
}

func addDagHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	decodedDag := &dag.Dag{}
	err = yaml.Unmarshal(body, decodedDag)

	if err != nil {
		http.Error(w, "Error decoding yaml", http.StatusInternalServerError)
		return
	}

	response := AddDagResponse{
		Message: "Successfully added DAG",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	dagList[decodedDag.Name] = decodedDag

	w.Write(jsonResponse)
}
