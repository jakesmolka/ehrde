/* (open) Electronic Health Record Data Explorer (or so)
--------------------------------------------------------
EHRDE source code as part of Jake Smolka's bachelor thesis

Following words are used frequently and with a certain meaning:
node := element to be displayed in widget, json or other data model or structure
template := openEHR data class
archetype := openEHR data class, templates aggregate archetypes
datavalues := openEHR/Think!EHR data_value class, archetypes aggregate datavalues
--------------------------------------------------------

main.go
v1: main.go contains the top level functions so the most basic program flow is shown.
*/

package main

import (
	d "./data"   //EHRDE data package
	s "./server" //EHRDE server package
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// prepares data to be displayed in widgets
func prepareWidgetData(dataMap map[string]interface{}) ([]interface{}, []interface{}, []interface{}) {
	//need to generate
	// treeWidgetSet &
	// visNodes, visEdges
	//and transform to:
	// widgetDataJson &
	// visNodesJson, visEdgesJson
	log.Print("########## Starting to build widget data.")

	numNodes := d.CountNodes(dataMap)
	treeWidgetSet, err := d.GenerateTreeData(dataMap, numNodes)
	errorHandler(err, "prepareWidgetData() -> generateTreeData()")
	log.Print("########## Done building tree widget data.")

	visNodes, visEdges, err := d.GenerateVisData(dataMap, numNodes)
	errorHandler(err, "prepareWidgetData() -> generateVisData()")
	log.Print("########## Done building network widget data.")

	return treeWidgetSet, visNodes, visEdges
}

// gets all data required by this app to minimize GETs and make later time-triggered refreshing possible
func getAllData() map[string]interface{} {
	log.Print("########## Staring to get all data.")
	data, err := d.HttpGetJsonThink("rest/v1/template")
	if err != nil {
		//
	}
	templatesSet := data.(map[string]interface{})["templates"] //gets array of map[string]interface{}'s

	// --- Getting and generating output for archetypes too
	dataMap, err := d.GetArchetypesByTemplates(templatesSet.([]interface{}))
	errorHandler(err, "getAllData() -> getArchetypesByTemplates()")
	log.Print("########## Done getting basic data. Continuing with AQL.")

	//adds fetching of *count* of each DATA_VALUE using AQL
	//NOTE: could be included into func above, but that way more modularity
	dataMap, err = d.AddCountValues(dataMap)
	errorHandler(err, "getAllData() -> getCountValues()")
	log.Print("########## Done with AQL. Writing json files to disk.")

	//NOTE: debug log-style ouput of dataMap as JSON file
	dataMapFile, err := json.Marshal(dataMap)
	errorHandler(err, "")
	err = ioutil.WriteFile("dataMap.json", dataMapFile, 0644) //NOTE debug
	errorHandler(err, "")
	log.Print("########## Wrote dataMap.json")

	//generating separate json with all datavalues using the same ids as in the widgets
	//	used to get access to data for detail-popup through index.html JS
	log.Print("########## (More AQL... Part of prototype feature)")
	datavaluesMap := d.GetDatavaluesMap(dataMap)
	datavaluesFile, err := json.Marshal(datavaluesMap)
	errorHandler(err, "")
	err = ioutil.WriteFile("assets/js/datavalues.json", datavaluesFile, 0644) //needed for ajax call by modals
	errorHandler(err, "")
	log.Print("########## Wrote assets/js/datavales.json")

	return dataMap
}

//general error handler to print error and optional more details
func errorHandler(err error, details string) {
	if err != nil {
		//log
		fmt.Println("ERROR: ", err)
		if details != "" {
			fmt.Println("details: ", details)
		}
	}
}

func main() {
	//pre: load config file into runtime
	err := d.LoadConfig()

	// gets and builds data model - so not every reload hits the server
	// (optional: time triggered call getAllData to rebuild data every now and then)
	// (optional: maybe even save data in database and just add new stuff each server start)
	dataMap := getAllData()

	treeWidgetSet, visNodes, visEdges := prepareWidgetData(dataMap)

	//bring tree widget data into real json to be displayed in jqxtree
	treeWidgetJson, err := json.Marshal(treeWidgetSet)
	errorHandler(err, "rootHandler() - Marshal(treeWidgetSet)")

	// building json data for displaying in visjs network
	visNodesJson, err := json.Marshal(visNodes)
	errorHandler(err, "rootHandler() - Marshal(visNodes)")
	visEdgesJson, err := json.Marshal(visEdges)
	errorHandler(err, "rootHandler() - marshal(visEdges)")

	// init web app handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//because widget data is static it gets loaded into rootHandler here
		s.RootHandler(w, r, treeWidgetJson, visNodesJson, visEdgesJson)
	})
	// init assets handler for jqwidgets' files
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	log.Print("++++++++++ All done. Starting server on :8080.")
	http.ListenAndServe(":8080", nil)
}
