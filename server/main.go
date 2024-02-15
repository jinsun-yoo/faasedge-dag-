package main

import (
	"encoding/json"
	"faasedge-dag/server/dag"
	"faasedge-dag/server/scheduler"
	"fmt"
	"io/ioutil"
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

	dagName := r.Header.Get("dag-name")
	if dagName == "" {
		http.Error(w, "dag-name header is missing", http.StatusBadRequest)
		return
	}

	if dag, ok := dagList[dagName]; ok {
		vdag := &scheduler.VDag {
			ClientId:      clientId,
			DagDefinition: *dag,
			PDagMap:       make(map[string][]*scheduler.PDagDeployment),
		}

		scheduler.Scheduler.ScheduleDag(vdag)

		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Message string `json:"message"`
		}{
			Message: "DAG scheduled successfully",
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	} else {
		http.Error(w, "DAG not found", http.StatusNotFound)
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
