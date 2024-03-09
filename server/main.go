package main

import (
	"encoding/json"
	"faasedge-dag/server/dag"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// Define a struct for your API response
type AddDagResponse struct {
	Message string `json:"message"`
}

var dagList = make(map[string]*dag.Dag)

func main() {
	http.HandleFunc("/dag/add", addDagHandler)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

// Handler function for the "/api/hello" endpoint
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
