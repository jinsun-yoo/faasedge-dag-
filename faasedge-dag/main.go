package main

import (
	"encoding/json"
	"faasedge-dag/server/dag"
	"faasedge-dag/server/location"
	"faasedge-dag/server/scheduler"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/yaml.v2"
)

type AddDagResponse struct {
	Message string `json:"message"`
}

var dagList = make(map[string]*dag.Dag)
var locationTracker = location.LocationTracker{}

func main() {

	file, err := ioutil.ReadFile("map.json")
	if err != nil {
		log.Fatalf("Unable to read map.json: %v", err)
	}

	var mapData struct {
		LowerLeft location.Coord `json:"lowerLeft"`
		UpperRight location.Coord `json:"upperRight"`
		AreasOfInterest []struct {
			LowerLeft  location.Coord `json:"lowerLeft"`
			UpperRight location.Coord `json:"upperRight"`
		} `json:"areasOfInterest"`
	}

	err = json.Unmarshal(file, &mapData)
	if err != nil {
		log.Fatalf("Unable to unmarshal map.json: %v", err)
	}

	fmt.Println(mapData.LowerLeft)

	for _, area := range mapData.AreasOfInterest {
		locationTracker.AreaOfInterestList = append(locationTracker.AreaOfInterestList, location.AreaOfInterest{
			LowerLeftCoord:  area.LowerLeft,
			UpperRightCoord: area.UpperRight,
			Clients:         make(map[int]bool),
			Clientlist:      []int{},
		})
	}

	http.HandleFunc("/dag/add", addDagHandler)
	http.HandleFunc("/dag/schedule", scheduleDagHandler)
	http.HandleFunc("/dag/invoke", invokeDag)
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

func invokeDag(w http.ResponseWriter, r *http.Request) {
	log.Println("dag invocation")
	clientIdRaw := r.Header.Get("Client-Id")
	if clientIdRaw == "" {
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

	_, ok := dagList[dagName]
	
	if !ok {
		log.Println("DAG ("+dagName+") not found", http.StatusNotFound)
		http.Error(w, "DAG ("+dagName+") not found", http.StatusNotFound)
		return
	}

	var requestData struct {
		Coordinate location.Coord `json:"coordinate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Println("Error decoding coordinate")
		http.Error(w, "Error decoding coordinate", http.StatusBadRequest)
		return
	}

	clientId, err := strconv.Atoi(clientIdRaw)

	if err != nil {
		log.Println("ClientId is not an integer")
		http.Error(w, "ClientId is not an integer", http.StatusBadRequest)
		return
	}

	locationTracker.RegisterLocation(clientId, requestData.Coordinate)

	aoi := locationTracker.LookupAreaOfInterest(requestData.Coordinate)
	log.Println("Area of interest: ", aoi)

	relatedVehicles := locationTracker.RelatedVehicles(clientId)
	log.Println("Related Vehicles: ", relatedVehicles)
}